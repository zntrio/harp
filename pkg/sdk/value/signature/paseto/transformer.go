// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package paseto

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"errors"
	"fmt"

	pasetov3 "zntr.io/paseto/v3"
	pasetov4 "zntr.io/paseto/v4"

	"zntr.io/harp/v2/pkg/sdk/types"
)

type pasetoTransformer struct {
	key interface{}
}

// -----------------------------------------------------------------------------

func (d *pasetoTransformer) To(_ context.Context, input []byte) ([]byte, error) {
	if types.IsNil(d.key) {
		return nil, fmt.Errorf("paseto: signer key must not be nil")
	}

	var (
		out string
		err error
	)

	switch sk := d.key.(type) {
	case ed25519.PrivateKey:
		out, err = pasetov4.Sign(input, sk, nil, nil)
	case *ecdsa.PrivateKey:
		out, err = pasetov3.Sign(input, sk, nil, nil)
	default:
		return nil, errors.New("paseto: key is not supported")
	}
	if err != nil {
		return nil, fmt.Errorf("paseto: unable so sign input: %w", err)
	}

	// No error
	return []byte(out), nil
}

func (d *pasetoTransformer) From(_ context.Context, input []byte) ([]byte, error) {
	var (
		payload []byte
		err     error
	)

	switch sk := d.key.(type) {
	case ed25519.PublicKey:
		payload, err = pasetov4.Verify(string(input), sk, nil, nil)
	case *ecdsa.PublicKey:
		payload, err = pasetov3.Verify(string(input), sk, nil, nil)
	default:
		return nil, errors.New("paseto: key is not supported")
	}
	if err != nil {
		return nil, fmt.Errorf("paseto: unable so sign input: %w", err)
	}

	// No error
	return payload, nil
}
