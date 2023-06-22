// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/dchest/uniuri"
	"zntr.io/harp/v2/build/fips"
)

// Build information. Populated at build-time.
var (
	Name      = "unknown"
	AppName   = "unknown"
	Version   = "unknown"
	Commit    = "unknown"
	Branch    = "unknown"
	GoVersion = "unknown"
	BuildTags = "unknown"
)

// NewInfo returns a build information object.
func NewInfo() Info {
	sdkVersion := getSDKVersion()
	return Info{
		Name:           Name,
		ComponentName:  AppName,
		Version:        Version,
		GitBranch:      Branch,
		GitCommit:      Commit,
		BuildTags:      BuildTags,
		GoVersion:      fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		BuildDeps:      depsFromBuildInfo(),
		HarpSdkVersion: sdkVersion,
	}
}

// Map provides the iterable version information.
type Info struct {
	Name           string     `json:"name"`
	ComponentName  string     `json:"component_name"`
	Version        string     `json:"version"`
	GitBranch      string     `json:"branch"`
	GitCommit      string     `json:"commit"`
	BuildTags      string     `json:"build_tags"`
	GoVersion      string     `json:"go"`
	BuildDeps      []buildDep `json:"build_deps"`
	HarpSdkVersion string     `json:"harp_sdk_version,omitempty"`
}

// Full returns full composed version string.
func (i *Info) String() string {
	if fips.Enabled() {
		return fmt.Sprintf("%s [%s:%s] (Go: %s, FIPS Mode, Flags: %s)", i.Version, i.GitBranch, i.GitCommit, i.GoVersion, i.BuildTags)
	}

	return fmt.Sprintf("%s [%s:%s] (Go: %s, Flags: %s)", i.Version, i.GitBranch, i.GitCommit, i.GoVersion, i.BuildTags)
}

// JSON returns json representation of build info.
func (i *Info) JSON() string {
	payload, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(payload)
}

// ID returns an instance id.
func ID() string {
	return uniuri.NewLen(64)
}

// -----------------------------------------------------------------------------

func getSDKVersion() string {
	// Extract build info
	deps, ok := debug.ReadBuildInfo()
	if !ok {
		return "unable to read deps"
	}

	// Look for harp dependency version
	var sdkVersion string
	for _, dep := range deps.Deps {
		if dep.Path == "zntr.io/harp/v2" {
			sdkVersion = dep.Version
		}
	}

	return sdkVersion
}

func depsFromBuildInfo() (deps []buildDep) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	for _, dep := range buildInfo.Deps {
		deps = append(deps, buildDep{dep})
	}

	return
}

type buildDep struct {
	*debug.Module
}

func (d buildDep) String() string {
	if d.Replace != nil {
		return fmt.Sprintf("%s@%s => %s@%s %s", d.Path, d.Version, d.Replace.Path, d.Replace.Version, d.Replace.Sum)
	}

	return fmt.Sprintf("%s@%s %s", d.Path, d.Version, d.Sum)
}

func (d buildDep) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }
