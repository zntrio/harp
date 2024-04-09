// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/go-jose/go-jose/v3"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/signature"
)

func init() {
	signature.Register("jws", Transformer)
}

// Transformer returns a JWS signature value transformer instance.
func Transformer(key string) (value.Transformer, error) {
	// Remove prefix
	key = strings.TrimPrefix(key, "jws:")

	// Decode key
	keyRaw, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("unable to decode transformer key: %w", err)
	}

	// Check JWK encoding
	var jwk jose.JSONWebKey
	if errJSON := json.Unmarshal(keyRaw, &jwk); errJSON != nil {
		return nil, fmt.Errorf("unable to decode the transformer key: %w", errJSON)
	}

	// Return transformer implementation
	return &jwsTransformer{
		key: jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(jwk.Algorithm),
			Key:       &jwk,
		},
	}, nil
}

// -----------------------------------------------------------------------------

type nonceSource struct{}

func (n *nonceSource) Nonce() (string, error) {
	return uniuri.NewLen(8), nil
}
