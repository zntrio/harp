// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"regexp"
	"runtime"

	semver "github.com/Masterminds/semver/v3"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/log"
)

var versionSemverRe = regexp.MustCompile("[0-9.]+")

// Is return true if current go version is included in given array.
func Is(constraints ...string) bool {
	// Extract version digit from go runtime version.
	v := versionSemverRe.FindString(runtime.Version())
	if v == "" {
		panic("unable to extract go runtime version")
	}

	// Parse golang version as semver
	sv := semver.MustParse(v)

	// Parse all constraints and check according to go version.
	for _, c := range constraints {
		constraint, err := semver.NewConstraint(c)
		if err != nil {
			log.Bg().Error("unable to parse version constraint", zap.String("constraint", c))
			return false
		}

		// Check version
		if constraint.Check(sv) {
			return true
		}
	}

	// No match found
	return false
}
