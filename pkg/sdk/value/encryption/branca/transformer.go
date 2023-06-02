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
