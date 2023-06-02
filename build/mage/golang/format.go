// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Format all source code.
func Format() error {
	mg.Deps(CollectGoFiles)

	color.Cyan("## Format everything")

	for pth := range CollectedGoFiles {
		args := []string{"-w"}
		args = append(args, pth)

		if err := sh.RunV("gofumpt", args...); err != nil {
			return err
		}
	}

	return nil
}
