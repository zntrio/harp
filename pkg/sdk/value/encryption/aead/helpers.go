// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package aead

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

const (
	keyLength = 32
)

func encrypt(ctx context.Context, plaintext []byte, ciph cipher.AEAD) ([]byte, error) {
	if len(plaintext) > 64*1024*1024 {
		return nil, errors.New("value too large")
	}
	nonce := make([]byte, ciph.NonceSize(), ciph.NonceSize()+ciph.Overhead()+len(plaintext))
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("unable to generate nonce: %w", err)
	}

	// Retrieve additional data from context
	aad, _ := encryption.AdditionalData(ctx)

	cipherText := ciph.Seal(nil, nonce, plaintext, aad)

	return append(nonce, cipherText...), nil
}

func decrypt(ctx context.Context, ciphertext []byte, ciph cipher.AEAD) ([]byte, error) {
	if len(ciphertext) < ciph.NonceSize() {
		return nil, errors.New("ciphered text too short")
	}

	nonce := ciphertext[:ciph.NonceSize()]
	text := ciphertext[ciph.NonceSize():]

	// Retrieve additional data from context
	aad, _ := encryption.AdditionalData(ctx)

	clearText, err := ciph.Open(nil, nonce, text, aad)
	if err != nil {
		return nil, errors.New("failed to decrypt given message")
	}

	return clearText, nil
}
