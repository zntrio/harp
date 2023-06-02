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

// Lint all source code.
func Lint(basePath string) func() error {
	return func() error {
		color.Cyan("## Lint go code")
		return sh.RunV("golangci-lint", "run", "-c", filepath.Clean(filepath.Join(basePath, ".golangci.yml")))
	}
}
