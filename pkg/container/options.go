// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import "github.com/awnumar/memguard"

// Option describes generate container operation options.
type Option func(opts *Options)

// Options defines the operation settings.
type Options struct {
	psk            *memguard.LockedBuffer
	peersPublicKey []string
}

// WithPreSharedKey sets the pre-sharey used for seal/unseal operations.
func WithPreSharedKey(psk *memguard.LockedBuffer) Option {
	return func(opts *Options) {
		opts.psk = psk
	}
}

// WithPeerPublicKeys sets the public key which are able to unseal the container.
func WithPeerPublicKeys(peers []string) Option {
	return func(opts *Options) {
		opts.peersPublicKey = peers
	}
}
