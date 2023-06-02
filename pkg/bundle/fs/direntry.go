// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import "io/fs"

type dirEntry struct {
	fi fs.FileInfo
}

// Compile time type assertion.
var _ fs.DirEntry = (*dirEntry)(nil)

// -----------------------------------------------------------------------------

func (d *dirEntry) Name() string {
	return d.fi.Name()
}

func (d *dirEntry) IsDir() bool {
	return d.fi.IsDir()
}

func (d *dirEntry) Type() fs.FileMode {
	return d.fi.Mode()
}

func (d *dirEntry) Info() (fs.FileInfo, error) {
	return d.fi, nil
}
