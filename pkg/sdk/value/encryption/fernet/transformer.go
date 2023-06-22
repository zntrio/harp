// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package fernet

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fernet/fernet-go"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

func init() {
	encryption.Register("fernet", Transformer)
}

// Transformer returns a fernet encryption transformer.
func Transformer(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "fernet:")

	// Check given keys
	k, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, fmt.Errorf("fernet: unable to initialize fernet transformer: %w", err)
	}

	// Return decorator constructor
	return &fernetTransformer{
		key: k,
	}, nil
}

// -----------------------------------------------------------------------------

type fernetTransformer struct {
	key *fernet.Key
}

func (d *fernetTransformer) To(_ context.Context, input []byte) ([]byte, error) {
	// Encrypt value
	out, err := fernet.EncryptAndSign(input, d.key)
	if err != nil {
		return nil, fmt.Errorf("fernet: unable to transform input value: %w", err)
	}

	// No error
	return out, nil
}

func (d *fernetTransformer) From(_ context.Context, input []byte) ([]byte, error) {
	// Encrypt value
	out := fernet.VerifyAndDecrypt(input, 0, []*fernet.Key{d.key})
	if out == nil {
		return nil, errors.New("fernet: unable to decrypt value")
	}

	// No error
	return out, nil
}
