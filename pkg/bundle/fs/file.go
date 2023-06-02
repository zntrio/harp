// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import (
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/awnumar/memguard"
)

type file struct {
	modTime    time.Time
	name       string
	bodyReader io.Reader
	size       int64
	content    *memguard.Enclave
	mode       os.FileMode
	closed     bool
}

// Compile time type assertion.
var _ fs.File = (*file)(nil)

// -----------------------------------------------------------------------------

func (f *file) Stat() (fs.FileInfo, error) {
	// Check file state
	if f.closed {
		return nil, fs.ErrClosed
	}

	// Return file information
	return &fileInfo{
		name:    f.name,
		size:    f.size,
		modTime: f.modTime,
		mode:    f.mode,
	}, nil
}

func (f *file) Read(b []byte) (int, error) {
	// Check file state
	if f.closed || f.bodyReader == nil {
		return 0, fs.ErrClosed
	}

	// Delegate to reader
	return f.bodyReader.Read(b)
}

func (f *file) Close() error {
	if f.closed {
		return fs.ErrClosed
	}
	f.closed = true
	f.bodyReader = nil
	return nil
}
