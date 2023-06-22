// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package identity

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"zntr.io/harp/v2/pkg/container/identity/key"
	"zntr.io/harp/v2/pkg/sdk/types"
)

const (
	apiVersion = "harp.zntr.io/v2"
	kind       = "ContainerIdentity"
)

// -----------------------------------------------------------------------------

type PrivateKeyGeneratorFunc func(io.Reader) (*key.JSONWebKey, string, error)

// New identity from description.
func New(random io.Reader, description string, generator PrivateKeyGeneratorFunc) (*Identity, []byte, error) {
	// Check arguments
	if err := validation.Validate(description, validation.Required, is.ASCII); err != nil {
		return nil, nil, fmt.Errorf("unable to create identity with invalid description: %w", err)
	}

	// Delegate to generator
	jwk, encodedPub, err := generator(random)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate identity private key: %w", err)
	}

	// Encode JWK as json
	payload, err := json.Marshal(jwk)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to serialize identity keypair: %w", err)
	}

	// Prepae identity object
	id := &Identity{
		APIVersion:  apiVersion,
		Kind:        kind,
		Timestamp:   time.Now().UTC(),
		Description: description,
		Public:      encodedPub,
	}

	// Encode to json for signature
	protected, err := json.Marshal(id)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to serialize identity for signature: %w", err)
	}

	// Sign the protected data
	sig, err := jwk.Sign(protected)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to sign protected data: %w", err)
	}

	// Auto-assign the signature
	id.Signature = sig

	// Return unsealed identity
	return id, payload, nil
}

// FromReader extract identity instance from reader.
func FromReader(r io.Reader) (*Identity, error) {
	// Check arguments
	if types.IsNil(r) {
		return nil, fmt.Errorf("unable to read nil reader")
	}

	// Convert input as a map
	var input Identity
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return nil, fmt.Errorf("unable to decode input JSON: %w", err)
	}

	// Check component
	if input.Private == nil {
		return nil, fmt.Errorf("invalid identity: missing private component")
	}

	// Validate self signature
	if errVerify := input.Verify(); errVerify != nil {
		return nil, fmt.Errorf("unable to verify identity: %w", errVerify)
	}

	// Return no error
	return &input, nil
}
