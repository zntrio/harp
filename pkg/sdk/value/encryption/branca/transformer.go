// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package branca

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/essentialkaos/branca"

	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

func init() {
	encryption.Register("branca", Transformer)
}

// Transformer returns a branca encryption transformer.
func Transformer(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "branca:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("branca: unable to decode key: %w", err)
	}

	// Check given keys
	codec, err := branca.NewBranca(k)
	if err != nil {
		return nil, fmt.Errorf("branca: unable to initialize the key: %w", err)
	}

	// Return decorator constructor
	return &brancaTransformer{
		codec: codec,
	}, nil
}

// -----------------------------------------------------------------------------

type brancaTransformer struct {
	codec *branca.Branca
}

func (d *brancaTransformer) To(_ context.Context, input []byte) ([]byte, error) {
	// Encrypt value
	out, err := d.codec.Encode(input)
	if err != nil {
		return nil, fmt.Errorf("branca: unable to transform input value: %w", err)
	}

	// No error
	return out, nil
}

func (d *brancaTransformer) From(_ context.Context, input []byte) ([]byte, error) {
	// Decrypt value without checking expiration
	out, err := d.codec.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("branca: unable to decrypt branca token: %w", err)
	}

	// No error
	return out.Payload(), nil
}
