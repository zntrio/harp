// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package seal

import (
	"io"

	"github.com/awnumar/memguard"
	containerv1 "zntr.io/harp/v2/api/gen/go/harp/container/v1"
)

// Streategy describes the sealing/unsealing contract.
type Strategy interface {
	// CenerateKey create an key pair used as container identifier.
	GenerateKey(...GenerateOption) (publicKey, privateKey string, err error)
	// Seal the given container using the implemented algorithm.
	Seal(io.Reader, *containerv1.Container, ...string) (*containerv1.Container, error)
	// Seal the given container using the implemented algorithm.
	SealWithPSK(io.Reader, *containerv1.Container, *memguard.LockedBuffer, ...string) (*containerv1.Container, error)
	// Unseal the given container using the given identity.
	Unseal(c *containerv1.Container, id *memguard.LockedBuffer) (*containerv1.Container, error)
	// UnsealWithPSK unseals the given container using the given identity and the gievn preshared key.
	UnsealWithPSK(c *containerv1.Container, id *memguard.LockedBuffer, psk *memguard.LockedBuffer) (*containerv1.Container, error)
}

// GenerateOptions represents container key generation options.
type GenerateOptions struct {
	DCKDMasterKey *memguard.LockedBuffer
	DCKDTarget    string
	RandomSource  io.Reader
}

// GenerateOption represents functional pattern builder for optional parameters.
type GenerateOption func(o *GenerateOptions)

// WithDeterministicKey enables deterministic container key generation.
func WithDeterministicKey(masterKey *memguard.LockedBuffer, target string) GenerateOption {
	return func(o *GenerateOptions) {
		o.DCKDMasterKey = masterKey
		o.DCKDTarget = target
	}
}

// WithRandom provides the random source for key generation.
func WithRandom(random io.Reader) GenerateOption {
	return func(o *GenerateOptions) {
		o.RandomSource = random
	}
}
