// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package key

import (
	"crypto/ed25519"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"
	"io"

	"zntr.io/harp/v2/pkg/sdk/security/crypto/deterministicecdsa"

	"golang.org/x/crypto/nacl/box"
)

func Legacy(random io.Reader) (*JSONWebKey, string, error) {
	// Generate X25519 keys as identity
	pub, priv, err := box.GenerateKey(random)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate identity keypair: %w", err)
	}

	// Wrap as JWK
	return &JSONWebKey{
		Kty: "OKP",
		Crv: "X25519",
		X:   base64.RawURLEncoding.EncodeToString(pub[:]),
		D:   base64.RawURLEncoding.EncodeToString(priv[:]),
	}, base64.RawURLEncoding.EncodeToString(pub[:]), err
}

func Ed25519(random io.Reader) (*JSONWebKey, string, error) {
	// Generate ed25519 keys as identity
	pub, priv, err := ed25519.GenerateKey(random)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate identity keypair: %w", err)
	}

	// Wrap as JWK
	return &JSONWebKey{
		Kty: "OKP",
		Crv: "Ed25519",
		X:   base64.RawURLEncoding.EncodeToString(pub[:]),
		D:   base64.RawURLEncoding.EncodeToString(priv[:]),
	}, fmt.Sprintf("v1.ipk.%s", base64.RawURLEncoding.EncodeToString(pub[:])), err
}

func P384(random io.Reader) (*JSONWebKey, string, error) {
	// Generate ecdsa P-384 keys as identity
	priv, err := deterministicecdsa.GenerateKey(elliptic.P384(), random)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate identity keypair: %w", err)
	}

	// Marshall as compressed point
	pub := elliptic.MarshalCompressed(priv.Curve, priv.PublicKey.X, priv.PublicKey.Y)

	// Wrap as JWK
	return &JSONWebKey{
		Kty: "EC",
		Crv: "P-384",
		X:   base64.RawURLEncoding.EncodeToString(priv.PublicKey.X.Bytes()),
		Y:   base64.RawURLEncoding.EncodeToString(priv.PublicKey.Y.Bytes()),
		D:   base64.RawURLEncoding.EncodeToString(priv.D.Bytes()),
	}, fmt.Sprintf("v2.ipk.%s", base64.RawURLEncoding.EncodeToString(pub)), err
}
