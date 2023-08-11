// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"zntr.io/harp/v2/build/mage/git"
)

type buildOpts struct {
	binaryName  string
	packageName string
	cgoEnabled  bool
	pieEnabled  bool
	goOS        string
	goArch      string
	goArm       string
}

// BuildOption is used to define function option pattern.
type BuildOption func(*buildOpts)

// -----------------------------------------------------------------------------

// WithCGO enables CGO compilation.
func WithCGO() BuildOption {
	return func(opts *buildOpts) {
		opts.cgoEnabled = true
	}
}

// WithPIE enables Position Independent Executable compilation.
func WithPIE() BuildOption {
	return func(opts *buildOpts) {
		opts.pieEnabled = true
	}
}

// GOOS sets the GOOS value during build.
func GOOS(value string) BuildOption {
	return func(opts *buildOpts) {
		opts.goOS = value
	}
}

// GOARCH sets the GOARCH value during build.
func GOARCH(value string) BuildOption {
	return func(opts *buildOpts) {
		opts.goArch = value
	}
}

// GOARM sets the GOARM value during build.
func GOARM(value string) BuildOption {
	return func(opts *buildOpts) {
		opts.goArm = value
	}
}

// -----------------------------------------------------------------------------

// Build the given binary using the given package.
//
//nolint:funlen // to refactor
func Build(name, packageName, version string, opts ...BuildOption) func() error {
	const (
		defaultCgoEnabled = false
		defaultGoOs       = runtime.GOOS
		defaultGoArch     = runtime.GOARCH
		defaultGoArm      = ""
	)

	// Default build options
	defaultOpts := &buildOpts{
		binaryName:  name,
		packageName: packageName,
		cgoEnabled:  defaultCgoEnabled,
		goOS:        defaultGoOs,
		goArch:      defaultGoArch,
		goArm:       defaultGoArm,
	}

	// Apply options
	for _, o := range opts {
		o(defaultOpts)
	}

	return func() error {
		// Retrieve git info first
		mg.SerialDeps(git.CollectInfo)

		// Generate artifact name
		artifactName := fmt.Sprintf("%s-%s-%s%s", name, defaultOpts.goOS, defaultOpts.goArch, defaultOpts.goArm)

		// Compilation flags
		compilationFlags := []string{}

		// Check if fips is enabled
		buildTags := "-tags=!fips"
		if os.Getenv("HARP_BUILD_FIPS_MODE") == "1" {
			artifactName = fmt.Sprintf("%s-fips", artifactName)
			compilationFlags = append(compilationFlags, "fips")
			buildTags = "-tags=fips,goexperiment.boringcrypto"
		}

		// Check if CGO is enabled
		if defaultOpts.cgoEnabled {
			artifactName = fmt.Sprintf("%s-cgo", artifactName)
			compilationFlags = append(compilationFlags, "cgo")
		}

		// Enable PIE if requested
		buildMode := "-buildmode=exe"
		if defaultOpts.pieEnabled {
			buildMode = "-buildmode=pie"
			artifactName = fmt.Sprintf("%s-pie", artifactName)
			compilationFlags = append(compilationFlags, "pie")
		}

		// Check compilation flags
		strCompilationFlags := "defaults"
		if len(compilationFlags) > 0 {
			strCompilationFlags = strings.Join(compilationFlags, ",")
		}

		// Inject version information
		varsSetByLinker := map[string]string{
			"zntr.io/harp/v2/build/version.Name":      name,
			"zntr.io/harp/v2/build/version.AppName":   packageName,
			"zntr.io/harp/v2/build/version.Version":   version,
			"zntr.io/harp/v2/build/version.Commit":    git.Revision,
			"zntr.io/harp/v2/build/version.Branch":    git.Branch,
			"zntr.io/harp/v2/build/version.BuildTags": strCompilationFlags,
		}
		linkerKeys := make([]string, 0, len(varsSetByLinker))
		for k := range varsSetByLinker {
			linkerKeys = append(linkerKeys, k)
		}
		sort.Strings(linkerKeys)

		var linkerArgs []string
		for _, k := range linkerKeys {
			linkerArgs = append(linkerArgs, "-X", fmt.Sprintf("'%s=%s'", k, varsSetByLinker[k]))
		}

		// Strip and remove DWARF
		linkerArgs = append(linkerArgs, "-s", "-w", "-buildid=")

		// Assemble ldflags
		ldflagsValue := strings.Join(linkerArgs, " ")

		// Build environment
		env := map[string]string{
			"GOOS":        defaultOpts.goOS,
			"GOARCH":      defaultOpts.goArch,
			"CGO_ENABLED": "0",
			// Enable experimental memory management via arenas
			"GOEXPERIMENT": "arenas",
		}
		if defaultOpts.cgoEnabled {
			env["CGO_ENABLED"] = "1"
		}
		if defaultOpts.goArm != "" {
			env["GOARM"] = defaultOpts.goArm
		}

		// Create output directory
		if errMkDir := os.Mkdir("bin", 0o700); errMkDir != nil {
			if !errors.Is(errMkDir, os.ErrExist) {
				return fmt.Errorf("unable to create output directory: %w", errMkDir)
			}
		}

		// Generate output filename
		filename := fmt.Sprintf("bin/%s", artifactName)
		if defaultOpts.goOS == "windows" {
			filename = fmt.Sprintf("%s.exe", filename)
		}

		fmt.Fprintf(os.Stdout, " > Generating SBOM %s [%s] [os:%s arch:%s%s flags:%v tag:%v]\n", defaultOpts.binaryName, defaultOpts.packageName, defaultOpts.goOS, defaultOpts.goArch, defaultOpts.goArm, strCompilationFlags, version)

		// Generate SBOM
		if err := sh.RunWith(env, "cyclonedx-gomod", "app", "-json", "-output", fmt.Sprintf("%s.sbom.json", filename), "-files", "-licenses", "-main", fmt.Sprintf("cmd/%s", defaultOpts.binaryName), "-packages"); err != nil {
			return fmt.Errorf("unable to generate SBOM for artifact: %w", err)
		}

		fmt.Fprintf(os.Stdout, " > Building %s [%s] [os:%s arch:%s%s flags:%v tag:%v]\n", defaultOpts.binaryName, defaultOpts.packageName, defaultOpts.goOS, defaultOpts.goArch, defaultOpts.goArm, strCompilationFlags, version)

		// Compile
		return sh.RunWith(env, "go", "build", buildMode, buildTags, "-trimpath", "-mod=vendor", "-buildvcs=false", "-ldflags", ldflagsValue, "-o", filename, packageName)
	}
}
