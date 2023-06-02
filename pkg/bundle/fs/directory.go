// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"
)

type directory struct {
	sync.RWMutex

	name     string
	perm     os.FileMode
	modTime  time.Time
	children map[string]interface{}
}

// Compile time type assertion.
var _ fs.ReadDirFile = (*directory)(nil)

// -----------------------------------------------------------------------------

func (d *directory) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    d.name,
		size:    1,
		modTime: d.modTime,
		mode:    d.perm | fs.ModeDir,
	}, nil
}

func (d *directory) Read(b []byte) (int, error) {
	return 0, errors.New("is a directory")
}

func (d *directory) Close() error {
	return nil
}

func (d *directory) ReadDir(n int) ([]fs.DirEntry, error) {
	// Lock for read
	d.RLock()
	defer d.RUnlock()

	// Retrieve children entry count
	childrenNames := []string{}
	for entryName := range d.children {
		childrenNames = append(childrenNames, entryName)
	}

	// Apply read limit
	if n <= 0 {
		n = len(childrenNames)
	}

	// Iterate on children entities
	out := []fs.DirEntry{}
	for i := 0; i < len(childrenNames) && i < n; i++ {
		name := childrenNames[i]
		h := d.children[name]

		switch item := h.(type) {
		case *directory:
			out = append(out, &dirEntry{
				fi: &fileInfo{
					name: item.name,
					size: 1,
					mode: item.perm | os.ModeDir,
				},
			})
		case *file:
			out = append(out, &dirEntry{
				fi: &fileInfo{
					name:    item.name,
					size:    item.size,
					modTime: item.modTime,
					mode:    item.mode,
				},
			})
		default:
			continue
		}
	}

	// Check directory entry exhaustion
	if n > len(childrenNames) {
		return out, io.EOF
	}

	// Check empty response
	if len(out) == 0 {
		return out, errors.New("directory has no entry")
	}

	// Return result
	return out, nil
}
