// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbox

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"zntr.io/harp/v2/build/fips"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

func init() {
	if !fips.Enabled() {
		encryption.Register("secretbox", Transformer)
	}
}

// Transformer returns a Nacl SecretBox encryption value transformer.
func Transformer(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "secretbox:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("secretbox: unable to decode key: %w", err)
	}
	if l := len(k); l != keyLength {
		return nil, fmt.Errorf("secretbox: invalid secret key length (%d)", l)
	}

	// Copy secret key
	secretKey := new([keyLength]byte)
	copy(secretKey[:], k)

	// Return transformer
	return &secretboxTransformer{
		key: secretKey,
	}, nil
}

// -----------------------------------------------------------------------------

type secretboxTransformer struct {
	key *[keyLength]byte
}

func (d *secretboxTransformer) From(_ context.Context, input []byte) ([]byte, error) {
	// Check output
	if l := len(input); l < nonceLength {
		return nil, fmt.Errorf("secretbox: invalid secret length (%d), check encryption status", l)
	}

	// Decrypt value
	out, err := decrypt(input, *d.key)
	if err != nil {
		return nil, fmt.Errorf("secretbox: unable to transform value: %w", err)
	}

	// No error
	return out, nil
}

func (d *secretboxTransformer) To(_ context.Context, input []byte) ([]byte, error) {
	// Encrypt value
	out, err := encrypt(input, *d.key)
	if err != nil {
		return nil, fmt.Errorf("secretbox: unable to transform value: %w", err)
	}

	// No error
	return out, nil
}
