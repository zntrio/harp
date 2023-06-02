// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package archive

type createOptions struct {
	rootPath     string
	includeGlobs []string
	excludeGlobs []string
}

type CreateOption func(opts *createOptions)

// WithCreateRootPath sets the root path for archive creation.
func WithCreateRootPath(value string) CreateOption {
	return func(opts *createOptions) {
		opts.rootPath = value
	}
}

// IncludeFiles sets the file inclusion filter values.
func IncludeFiles(values ...string) CreateOption {
	return func(opts *createOptions) {
		opts.includeGlobs = values
	}
}

// ExcludeFiles sets the file exclusion filter values.
func ExcludeFiles(values ...string) CreateOption {
	return func(opts *createOptions) {
		opts.includeGlobs = values
	}
}
