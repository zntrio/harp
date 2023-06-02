// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package raw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/square/go-jose.v2"

	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/signature"
)

func init() {
	signature.Register("raw", Transformer)
}

// Transformer returns a JWS signature value transformer instance.
func Transformer(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "raw:")

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

	// Delegate to transformer
	return &rawTransformer{
		key: jwk.Key,
	}, err
}
