// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"zntr.io/harp/v2/build/artifact"
	"zntr.io/harp/v2/build/mage/docker"
	"zntr.io/harp/v2/build/mage/git"
	"zntr.io/harp/v2/build/mage/golang"
	"zntr.io/harp/v2/build/mage/release"
)

// -----------------------------------------------------------------------------

type Code mg.Namespace

// Lint code using golangci-lint.
func (Code) Lint() {
	mg.Deps(Code.Format)

	color.Red("## Lint source")
	mg.Deps(golang.Lint("."))
}

// Format source code and process imports.
func (Code) Format() {
	color.Red("## Formatting all sources")
	mg.SerialDeps(golang.Format, golang.Import)
}

// Licenser apply copyright banner to source code.
func (Code) Licenser() error {
	mg.SerialDeps(golang.Format, golang.Import)

	color.Red("## Add license banner")
	return sh.RunV("go-licenser")
}

// Generate SDK code (mocks, tests, etc.)
func (Code) Generate() {
	color.Cyan("## Generate code")
	mg.SerialDeps(
		func() error {
			return golang.Generate("SDK", "zntr.io/harp/v2/pkg/...")()
		},
	)
}

// -----------------------------------------------------------------------------

type API mg.Namespace

// Generate protobuf objects from proto definitions.
func (API) Generate() error {
	color.Blue("### Regenerate API")
	if err := sh.RunV("task", "-d", "api"); err != nil {
		return err
	}

	mg.SerialDeps(Code.Licenser)
	return nil
}

// -----------------------------------------------------------------------------

var Default = Build

var (
	harpCli = &artifact.Command{
		Package:     "zntr.io/harp/v2",
		Name:        "Harp",
		Description: "Secret management toolchain",
	}
)

// Build harp executable.
func Build() error {
	banner := figure.NewFigure("Harpocrates", "", true)
	banner.Print()

	fmt.Println("")
	color.Red("# Build Info ---------------------------------------------------------------")
	fmt.Printf("Go version : %s\n", runtime.Version())

	version, err := git.TagMatch("v*")
	if err != nil {
		return err
	}

	fmt.Printf("Git tag    : %s\n", version)

	color.Red("# Pipeline -----------------------------------------------------------------")
	mg.SerialDeps(golang.Vendor, golang.License("."), Code.Generate, golang.Lint("."), Test.Unit)

	color.Red("# Artifact(s) --------------------------------------------------------------")
	mg.Deps(Compile)

	// No error
	return nil
}

type Test mg.Namespace

// Test harp application.
func (Test) Unit() {
	color.Cyan("## Unit Tests")
	mg.SerialDeps(
		func() error {
			return golang.UnitTest("zntr.io/harp/v2/pkg/...")()
		},
		func() error {
			return golang.UnitTest("zntr.io/harp/v2/cmd/harp/...")()
		},
	)
}

// Test harp application.
func (Test) CLI() {
	color.Cyan("## CLI Tests")
	mg.SerialDeps(
		func() error {
			return golang.UnitTest("zntr.io/harp/v2/test/cmd")()
		},
	)
}

// Compile harp code to create an executable.
func Compile() error {
	// Extract git version
	version, err := git.TagMatch("v*")
	if err != nil {
		return err
	}

	// Build artifact
	return golang.Build("harp", "zntr.io/harp/v2/cmd/harp", version)()
}

// Release harp version and cross-compile code to produce all artifacts.
// RELEASE environment variable must be set to matching git tag.
func Release(ctx context.Context) error {
	color.Red(fmt.Sprintf("# Releasing (%s) ---------------------------------------------------", os.Getenv("RELEASE")))

	// Extract git version
	version, err := git.TagMatch("v*")
	if err != nil {
		return err
	}

	color.Cyan("## Cross compiling artifact")

	mg.CtxDeps(ctx,
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("darwin"), golang.GOARCH("amd64"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("darwin"), golang.GOARCH("arm64"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("linux"), golang.GOARCH("amd64"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("linux"), golang.GOARCH("arm"), golang.GOARM("7"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("linux"), golang.GOARCH("arm64"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("windows"), golang.GOARCH("amd64"),
			)()
		},
		func() error {
			return golang.Release(
				"harp",
				"zntr.io/harp/v2/cmd/harp",
				version,
				golang.GOOS("windows"), golang.GOARCH("arm64"),
			)()
		},
	)

	return ctx.Err()
}

// Homebrew generates homebrew formula from compiled artifacts.
func Homebrew() error {
	return release.HomebrewFormula(harpCli)()
}

// -----------------------------------------------------------------------------

type Docker mg.Namespace

// Tools prepares docker images with go toolchain and project tools.
func (Docker) Tools() error {
	return docker.Tools()
}

// Harp build harp docker image
func (Docker) Harp() error {
	return docker.Build(harpCli)()
}

// -----------------------------------------------------------------------------

type Releaser mg.Namespace

// Harp releases harp artifacts using docker pipeline.
func (Releaser) Harp() error {
	return docker.Release(harpCli)()
}
