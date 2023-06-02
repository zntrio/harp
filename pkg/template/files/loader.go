// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package files

import (
	"fmt"
	"io/fs"
)

// BufferedFile represents an archive file buffered for later processing.
type BufferedFile struct {
	Name string
	Data []byte
}

// ContentLoader loads file content.
type ContentLoader interface {
	Load() ([]*BufferedFile, error)
}

// Loader returns a new BufferedFile list from given path name.
func Loader(filesystem fs.FS, name string) (ContentLoader, error) {
	// Check if it's a directory
	fi, err := fs.Stat(filesystem, name)
	if err != nil {
		return nil, fmt.Errorf("unable to get file info for %q: %w", name, err)
	}

	// Is directory
	if fi.IsDir() {
		return &DirLoader{
			filesystem: filesystem,
			name:       name,
		}, nil
	}

	return nil, fmt.Errorf("only directory is supported as content loader")
}
