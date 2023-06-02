// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package files

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// DirLoader loads a chart from a directory.
type DirLoader struct {
	filesystem fs.FS
	name       string
}

// Load loads the chart.
func (l DirLoader) Load() ([]*BufferedFile, error) {
	return LoadDir(l.filesystem, l.name)
}

// LoadDir loads from a directory.
//
// This loads charts only from directories.
func LoadDir(filesystem fs.FS, dir string) ([]*BufferedFile, error) {
	// Check if path is valid
	if !fs.ValidPath(dir) {
		return nil, fmt.Errorf("%q is not a valid path", dir)
	}

	result := []*BufferedFile{}
	topdir := dir

	walk := func(name string, d fs.DirEntry, errWalk error) error {
		// Check walk error
		if errWalk != nil {
			return errWalk
		}

		// Compute relative path
		n, err := filepath.Rel(topdir, name)
		if err != nil {
			return fmt.Errorf("unable to compute relative path: %w", err)
		}
		if n == "" {
			return nil
		}

		// Normalize filepath
		n = filepath.ToSlash(n)

		// Ignore if it is a directory
		if d.IsDir() {
			return nil
		}

		// Irregular files include devices, sockets, and other uses of files that
		// are not regular files.
		if !d.Type().IsRegular() {
			return fmt.Errorf("cannot load irregular file %s as it has file mode type bits set", name)
		}

		// Read file content
		data, err := fs.ReadFile(filesystem, name)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", name, err)
		}

		// Append to result
		result = append(result, &BufferedFile{Name: n, Data: data})

		// No error
		return nil
	}
	if err := fs.WalkDir(filesystem, topdir, walk); err != nil {
		return nil, fmt.Errorf("unable to walk directory %q : %w", topdir, err)
	}

	// No error
	return result, nil
}
