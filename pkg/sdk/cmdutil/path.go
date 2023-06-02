// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrNoHome is raised when tilde expansion failed.
var ErrNoHome = errors.New("no home found")

// Expand a given path using `~` notation for HOMEDIR.
func Expand(path string) (string, error) {
	// Check condition
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	// Retrieve HOMEDIR
	home, err := getHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user home directory path: %w", err)
	}

	// Return result
	return home + path[1:], nil
}

func getHomeDir() (string, error) {
	var home string

	switch runtime.GOOS {
	case "windows":
		// Retrieve windows specific env
		home = filepath.Join(os.Getenv("HomeDrive"), os.Getenv("HomePath"))
		if home == "" {
			home = os.Getenv("UserProfile")
		}

	default:
		home = os.Getenv("HOME")
	}

	// Homedir not evaluable ?
	if home == "" {
		return "", ErrNoHome
	}

	// Return result
	return home, nil
}
