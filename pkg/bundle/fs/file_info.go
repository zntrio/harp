// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import (
	"io/fs"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    fs.FileMode
}

// Compile time type assertion.
var _ fs.FileInfo = (*fileInfo)(nil)

// -----------------------------------------------------------------------------

func (fi *fileInfo) Name() string {
	return fi.name
}

func (fi *fileInfo) Size() int64 {
	return fi.size
}

func (fi *fileInfo) Mode() fs.FileMode {
	return fi.mode
}

func (fi *fileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *fileInfo) IsDir() bool {
	return fi.mode&fs.ModeDir > 0
}

func (fi *fileInfo) Sys() interface{} {
	return nil
}
