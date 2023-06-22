// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jwe

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/square/go-jose.v2"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

type KeyAlgorithm string

//nolint:revive,stylecheck // Accepted case
var (
	AES128_KW          KeyAlgorithm = "a128kw"
	AES192_KW          KeyAlgorithm = "a192kw"
	AES256_KW          KeyAlgorithm = "a256kw"
	PBES2_HS256_A128KW KeyAlgorithm = "pbes2-hs256-a128kw"
	PBES2_HS384_A192KW KeyAlgorithm = "pbes2-hs384-a192kw"
	PBES2_HS512_A256KW KeyAlgorithm = "pbes2-hs512-a256kw"
)

func init() {
	encryption.Register("jwe", FromKey)
}

// FromKey returns an encryption transformer instance according to the given key format.
func FromKey(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "jwe:")

	switch {
	case strings.HasPrefix(key, "a128kw:"):
		return Transformer(AES128_KW, strings.TrimPrefix(key, "a128kw:"))
	case strings.HasPrefix(key, "a192kw:"):
		return Transformer(AES192_KW, strings.TrimPrefix(key, "a192kw:"))
	case strings.HasPrefix(key, "a256kw:"):
		return Transformer(AES256_KW, strings.TrimPrefix(key, "a256kw:"))
	case strings.HasPrefix(key, "pbes2-hs256-a128kw:"):
		return Transformer(PBES2_HS256_A128KW, strings.TrimPrefix(key, "pbes2-hs256-a128kw:"))
	case strings.HasPrefix(key, "pbes2-hs384-a192kw:"):
		return Transformer(PBES2_HS384_A192KW, strings.TrimPrefix(key, "pbes2-hs384-a192kw:"))
	case strings.HasPrefix(key, "pbes2-hs512-a256kw:"):
		return Transformer(PBES2_HS512_A256KW, strings.TrimPrefix(key, "pbes2-hs512-a256kw:"))
	default:
	}

	// Unsupported encryption scheme.
	return nil, fmt.Errorf("unsupported jwe algorithm for key %q", key)
}

// TransformerKey assemble a transformer key.
func TransformerKey(algorithm KeyAlgorithm, key string) string {
	return fmt.Sprintf("%s:%s", algorithm, key)
}

// Transformer returns a JWE encryption value transformer instance.
func Transformer(algorithm KeyAlgorithm, key string) (value.Transformer, error) {
	switch algorithm {
	case AES128_KW:
		// Try to decode the key
		k, err := base64.URLEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("jwe: unable to decode key: %w", err)
		}
		if len(k) < 16 {
			return nil, errors.New("jwe: key too short")
		}
		return transformer(k, jose.A128KW, jose.A128GCM)
	case AES192_KW:
		// Try to decode the key
		k, err := base64.URLEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("jwe: unable to decode key: %w", err)
		}
		if len(k) < 24 {
			return nil, errors.New("jwe: key too short")
		}
		return transformer(k, jose.A192KW, jose.A192GCM)
	case AES256_KW:
		// Try to decode the key
		k, err := base64.URLEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("jwe: unable to decode key: %w", err)
		}
		if len(k) < 32 {
			return nil, errors.New("jwe: key too short")
		}
		return transformer(k, jose.A256KW, jose.A256GCM)
	case PBES2_HS256_A128KW:
		return transformer(key, jose.PBES2_HS256_A128KW, jose.A128GCM)
	case PBES2_HS384_A192KW:
		return transformer(key, jose.PBES2_HS384_A192KW, jose.A192GCM)
	case PBES2_HS512_A256KW:
		return transformer(key, jose.PBES2_HS512_A256KW, jose.A256GCM)
	default:
	}

	// Unsupported encryption scheme.
	return nil, fmt.Errorf("unsupported jwe algorithm %q", algorithm)
}
