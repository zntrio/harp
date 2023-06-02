// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import "io/fs"

// BundleFS describes bundle filesystem contract.
type BundleFS interface {
	fs.ReadFileFS
	fs.ReadDirFS
	fs.StatFS
}
