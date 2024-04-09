// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jws

import (
	"context"
	"fmt"

	"github.com/go-jose/go-jose/v3"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/sdk/value/signature"
)

type jwsTransformer struct {
	key jose.SigningKey
}

// -----------------------------------------------------------------------------

func (d *jwsTransformer) To(ctx context.Context, input []byte) ([]byte, error) {
	if types.IsNil(d.key.Key) {
		return nil, fmt.Errorf("jws: signer key must not be nil")
	}

	opts := &jose.SignerOptions{}

	// If not deterministic add nonce in the protected header
	if !signature.IsDeterministic(ctx) {
		opts.NonceSource = &nonceSource{}
	}

	// Initialize a signer
	signer, err := jose.NewSigner(d.key, opts)
	if err != nil {
		return nil, fmt.Errorf("jws: unable to initialize a signer: %w", err)
	}

	// Sign input
	sig, err := signer.Sign(input)
	if err != nil {
		return nil, fmt.Errorf("jws: unable to sign the content: %w", err)
	}

	// Serialize content
	out, errSerialization := sig.CompactSerialize()

	if errSerialization != nil {
		return nil, fmt.Errorf("jws: unable to serialize final payload: %w", errSerialization)
	}

	// No error
	return []byte(out), nil
}

func (d *jwsTransformer) From(ctx context.Context, input []byte) ([]byte, error) {
	// Parse the signed object
	sig, err := jose.ParseSigned(string(input))
	if err != nil {
		return nil, fmt.Errorf("jws: unable to parse input: %w", err)
	}

	// Verify signature
	payload, err := sig.Verify(d.key.Key)
	if err != nil {
		return nil, fmt.Errorf("jws: unable to validate signature: %w", err)
	}

	// No error
	return payload, nil
}
