// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// Generate invoke the go:generate task on given package.
func Generate(name, packageName string) func() error {
	return func() error {
		fmt.Fprintf(os.Stdout, " > %s [%s]\n", name, packageName)
		return sh.RunV("go", "generate", packageName)
	}
}
