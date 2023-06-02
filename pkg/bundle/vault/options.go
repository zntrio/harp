// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"fmt"
	"regexp"
)

type options struct {
	prefix             string
	withSecretMetadata bool
	withVaultMetadata  bool
	workerCount        int64
	exclusions         []*regexp.Regexp
	includes           []*regexp.Regexp
	continueOnError    bool
}

// Option defines the functional pattern for bundle operation settings.
type Option func(*options) error

// WithExcludePath register a path exclusion regexp.
func WithExcludePath(value string) Option {
	return func(opts *options) error {
		// Compile RegExp first
		r, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("unable to compile `%s` as a valid regexp: %w", value, err)
		}

		// Append to exclusions
		opts.exclusions = append(opts.exclusions, r)

		// No error
		return nil
	}
}

// WithIncludePath register a path inclusion regexp.
func WithIncludePath(value string) Option {
	return func(opts *options) error {
		// Compile RegExp first
		r, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("unable to compile `%s` as a valid regexp: %w", value, err)
		}

		// Append to exclusions
		opts.includes = append(opts.includes, r)

		// No error
		return nil
	}
}

// WithPrefix add a prefix to path value.
func WithPrefix(value string) Option {
	return func(opts *options) error {
		opts.prefix = value
		// No error
		return nil
	}
}

// WithSecretMetadata add package metadata as secret value to be exported in Vault.
func WithSecretMetadata(value bool) Option {
	return func(opts *options) error {
		opts.withSecretMetadata = value
		// No error
		return nil
	}
}

// WithVaultMetadata add package metadata as secret metadata to be exported in Vault.
func WithVaultMetadata(value bool) Option {
	return func(opts *options) error {
		opts.withVaultMetadata = value
		// No error
		return nil
	}
}

// WithMaxWorkerCount sets the maximum count of active operation worker count.
func WithMaxWorkerCount(value int64) Option {
	return func(opts *options) error {
		opts.workerCount = value
		// No error
		return nil
	}
}

// WithContinueOnError enable/disbale stop processing on failure.
func WithContinueOnError(value bool) Option {
	return func(opts *options) error {
		opts.continueOnError = value
		// No error
		return nil
	}
}
