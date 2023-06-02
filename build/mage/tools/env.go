// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package tools

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// Env sets the environment for tools.
func Env() error {
	// Get current working directory
	name, err := os.Getwd()
	if err != nil {
		return err
	}

	// Get absolute path
	p, err := filepath.Abs(path.Join(name, "tools", "bin"))
	if err != nil {
		return err
	}

	// Add local bin in PATH
	return os.Setenv("PATH", fmt.Sprintf("%s:%s", p, os.Getenv("PATH")))
}
