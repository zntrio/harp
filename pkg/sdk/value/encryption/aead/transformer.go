// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package aead

import (
	"context"
	"crypto/cipher"
)

// -----------------------------------------------------------------------------

type aeadTransformer struct {
	aead cipher.AEAD
}

func (t *aeadTransformer) To(ctx context.Context, input []byte) ([]byte, error) {
	// Encrypt
	out, err := encrypt(ctx, input, t.aead)
	if err != nil {
		return nil, err
	}

	// Return result
	return out, nil
}

func (t *aeadTransformer) From(ctx context.Context, input []byte) ([]byte, error) {
	// Decrypt
	out, err := decrypt(ctx, input, t.aead)
	if err != nil {
		return nil, err
	}

	// No error
	return out, nil
}
