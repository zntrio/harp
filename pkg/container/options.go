// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package container

import "github.com/awnumar/memguard"

// Option describes generate container operation options.
type Option func(opts *Options)

// Options defines the operation settings
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

// WithPeerPublicKeys sets the public key which are able to unseal the container
func WithPeerPublicKeys(peers []string) Option {
	return func(opts *Options) {
		opts.peersPublicKey = peers
	}
}
