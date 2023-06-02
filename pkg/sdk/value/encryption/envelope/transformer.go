// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package envelope

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"

	"golang.org/x/crypto/cryptobyte"
)

// Transformer returns an envelope encryption value transformer.
func Transformer(envelopeService Service, transformerFactory encryption.TransformerFactoryFunc) (value.Transformer, error) {
	return &envelopeTransformer{
		envelopeService:        envelopeService,
		transformerFactoryFunc: transformerFactory,
	}, nil
}

// -----------------------------------------------------------------------------

type envelopeTransformer struct {
	envelopeService        Service
	transformerFactoryFunc encryption.TransformerFactoryFunc
}

func (t *envelopeTransformer) To(ctx context.Context, input []byte) ([]byte, error) {
	// Generate a random 32 byte length key
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return nil, fmt.Errorf("envelope: unable to generate dek key: %w", err)
	}

	// Encrypt DEK with envelope service
	encKey, err := t.envelopeService.Encrypt(ctx, newKey)
	if err != nil {
		return nil, fmt.Errorf("envelope: unable to encrypt dek: %w", err)
	}

	// Build a transformer using key
	transformer, err := t.transformerFactoryFunc(base64.URLEncoding.EncodeToString(newKey))
	if err != nil {
		return nil, fmt.Errorf("envelope: unable to initialize payload transformer: %w", err)
	}

	// Encrypt input using DEK
	payload, err := transformer.To(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("envelope: unable to encrypt data using dek: %w", err)
	}

	// Append the length of the encrypted DEK as the first 2 bytes.
	b := cryptobyte.NewBuilder(nil)
	b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
		b.AddBytes(encKey)
	})
	b.AddBytes(payload)

	// No error
	return b.Bytes()
}

func (t *envelopeTransformer) From(ctx context.Context, input []byte) ([]byte, error) {
	// Extract encrypted Data Encryption Key from input
	var encKey cryptobyte.String

	s := cryptobyte.String(input)
	if ok := s.ReadUint16LengthPrefixed(&encKey); !ok {
		return nil, fmt.Errorf("envelope: unable to read prefix")
	}

	// Encoded payload
	payload := []byte(s)

	// Decrypt DEK with envelope service
	key, err := t.envelopeService.Decrypt(ctx, encKey)
	if err != nil {
		return nil, fmt.Errorf("envelope: unable to decrypt dek: %w", err)
	}

	// Build a transformer using decoded key
	transformer, err := t.transformerFactoryFunc(base64.URLEncoding.EncodeToString(key))
	if err != nil {
		return nil, fmt.Errorf("envelope: unable to initialize payload transformer: %w", err)
	}

	// Delegate to transformer
	return transformer.From(ctx, payload)
}
