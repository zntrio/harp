// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// License checks allowed license of vendored dependencies.
func License(basePath string) func() error {
	return func() error {
		color.Cyan("## Check license")
		return sh.RunV("wwhrd", "check", "-f", filepath.Clean(filepath.Join(basePath, ".wwhrd.yml")))
	}
}
