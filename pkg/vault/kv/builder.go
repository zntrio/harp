// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"

	vpath "zntr.io/harp/v2/pkg/vault/path"
)

// Option defines the functional option pattern.
type Option func(opts *Options)

// Options defiens the default option value.
type Options struct {
	useCustomMetadata bool
	ctx               context.Context
}

// WithContext adds given context to all queries.
func WithContext(ctx context.Context) Option {
	return func(opts *Options) {
		opts.ctx = ctx
	}
}

// WithVaultMetatadata enable/disable the custom metadata storage strategy (requires Vault >=1.9).
func WithVaultMetatadata(value bool) Option {
	return func(opts *Options) {
		opts.useCustomMetadata = value
	}
}

// New build a KV service according to mountPath version.
func New(client *api.Client, path string, opts ...Option) (Service, error) {
	// Sanitize path
	secretPath := vpath.SanitizePath(path)

	// Defines default flag.
	dopts := &Options{
		useCustomMetadata: false,
		ctx:               context.Background(),
	}

	// Apply option function.
	for _, o := range opts {
		o(dopts)
	}

	// Detect mount path
	mountPath, v2, err := isKVv2(dopts.ctx, secretPath, client)
	if err != nil {
		return nil, fmt.Errorf("vault: unable to detect k/v backend version: %w", err)
	}

	// Build the service according to mountPath version
	var s Service
	if v2 {
		s = V2(client.Logical(), mountPath, dopts.useCustomMetadata)
	} else {
		s = V1(client.Logical(), mountPath)
	}

	// No error
	return s, nil
}
