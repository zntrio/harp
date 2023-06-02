// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"

	"zntr.io/harp/v2/build/mage/git"
)

// -----------------------------------------------------------------------------

// Release build and generate a final release artifact.
func Release(name, packageName, version string, opts ...BuildOption) func() error {
	return func() error {
		mg.SerialDeps(git.CollectInfo)

		// Retrieve release from ENV
		releaseVersion := os.Getenv("RELEASE")
		if releaseVersion == "" {
			return fmt.Errorf("RELEASE environment variable is missing")
		}

		// Release must be done on main branch only
		if git.Branch != "main" && os.Getenv("RELEASE_FORCE") == "" {
			return fmt.Errorf("a release must be build on 'main' branch only")
		}

		// Build the artifact
		if err := Build(
			name,
			packageName,
			version,
			opts...,
		)(); err != nil {
			return err
		}

		// No error
		return nil
	}
}
