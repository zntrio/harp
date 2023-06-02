// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package dae

import (
	"context"
	"crypto/cipher"
	"errors"
	"fmt"

	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

// -----------------------------------------------------------------------------

type daeTransformer struct {
	aead             cipher.AEAD
	nonceDeriverFunc NonceDeriverFunc
}

func (t *daeTransformer) To(ctx context.Context, input []byte) ([]byte, error) {
	// Check input size
	if len(input) > 64*1024*1024 {
		return nil, errors.New("value too large")
	}

	// Derive nonce
	nonce, err := t.nonceDeriverFunc(input, t.aead)
	if err != nil {
		return nil, fmt.Errorf("dae: unable to derive nonce: %w", err)
	}
	if len(nonce) != t.aead.NonceSize() {
		return nil, errors.New("dae: derived nonce is too short")
	}

	// Retrieve additional data from context
	aad, _ := encryption.AdditionalData(ctx)

	// Seal the cleartext with deterministic nonce
	cipherText := t.aead.Seal(nil, nonce, input, aad)

	// Return encrypted value
	return append(nonce, cipherText...), nil
}

func (t *daeTransformer) From(ctx context.Context, input []byte) ([]byte, error) {
	// Check input size
	if len(input) < t.aead.NonceSize() {
		return nil, errors.New("dae: ciphered text too short")
	}

	nonce := input[:t.aead.NonceSize()]
	text := input[t.aead.NonceSize():]
	aad, _ := encryption.AdditionalData(ctx)

	clearText, err := t.aead.Open(nil, nonce, text, aad)
	if err != nil {
		return nil, errors.New("failed to decrypt given message")
	}

	// No error
	return clearText, nil
}
