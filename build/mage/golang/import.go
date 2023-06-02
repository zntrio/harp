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

// Import fix all source code imports.
func Import() error {
	mg.Deps(CollectGoFiles)

	color.Cyan("## Process imports")

	for pth := range CollectedGoFiles {
		args := []string{"write", "-s", "Standard", "-s", "Prefix(golang.org/x/)", "-s", "Default", "-s", "Prefix(github.com/zntrio)"}
		args = append(args, pth)

		if err := sh.RunV("gci", args...); err != nil {
			return err
		}
	}

	return nil
}
