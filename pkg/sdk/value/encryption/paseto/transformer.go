// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package paseto

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"zntr.io/harp/v2/build/fips"
	pasetov4 "zntr.io/harp/v2/pkg/sdk/security/crypto/paseto/v4"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

func init() {
	if !fips.Enabled() {
		encryption.Register("paseto", Transformer)
	}
}

func Transformer(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "paseto:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("paseto: unable to decode key: %w", err)
	}
	if l := len(k); l != pasetov4.KeyLength {
		return nil, fmt.Errorf("paseto: invalid secret key length (%d)", l)
	}

	// Copy secret key
	var secretKey [pasetov4.KeyLength]byte
	copy(secretKey[:], k)

	return &pasetoTransformer{
		key: secretKey,
	}, nil
}

// -----------------------------------------------------------------------------

type pasetoTransformer struct {
	key [pasetov4.KeyLength]byte
}

func (d *pasetoTransformer) From(_ context.Context, input []byte) ([]byte, error) {
	return pasetov4.Decrypt(d.key[:], input, "", "")
}

func (d *pasetoTransformer) To(_ context.Context, input []byte) ([]byte, error) {
	// Encrypt with paseto v4.local
	return pasetov4.Encrypt(rand.Reader, d.key[:], input, "", "")
}
