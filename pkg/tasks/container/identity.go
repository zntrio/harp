// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"zntr.io/harp/v2/build/fips"
	"zntr.io/harp/v2/pkg/container/identity"
	"zntr.io/harp/v2/pkg/container/identity/key"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/tasks"
)

type IdentityVersion uint

const (
	LegacyIdentity IdentityVersion = 1
	ModernIdentity IdentityVersion = 2
	NISTIdentity   IdentityVersion = 3
)

// IdentityTask implements secret container identity creation task.
type IdentityTask struct {
	OutputWriter tasks.WriterProvider
	Description  string
	Transformer  value.Transformer
	Version      IdentityVersion
}

// Run the task.
func (t *IdentityTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if types.IsNil(t.Transformer) {
		return errors.New("unable to run task with a nil transformer")
	}
	if t.Description == "" {
		return fmt.Errorf("description must not be blank")
	}

	// Select appropriate strategy.
	var generator identity.PrivateKeyGeneratorFunc

	if fips.Enabled() {
		generator = key.P384
	} else {
		switch t.Version {
		case LegacyIdentity:
			generator = key.Legacy
		case ModernIdentity:
			generator = key.Ed25519
		case NISTIdentity:
			generator = key.P384
		default:
			return fmt.Errorf("invalid or unsupported identity version '%d'", t.Version)
		}
	}

	// Create identity
	id, payload, err := identity.New(rand.Reader, t.Description, generator)
	if err != nil {
		return fmt.Errorf("unable to create a new identity: %w", err)
	}

	// Encrypt the private key.
	identityPrivate, err := t.Transformer.To(ctx, payload)
	if err != nil {
		return fmt.Errorf("unable to encrypt the private identity key: %w", err)
	}

	// Assign private key
	id.Private = &identity.PrivateKey{
		Content: base64.RawURLEncoding.EncodeToString(identityPrivate),
	}

	// Retrieve output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve output writer handle: %w", err)
	}

	// Create identity output
	if err := json.NewEncoder(writer).Encode(id); err != nil {
		return fmt.Errorf("unable to serialize final identity: %w", err)
	}

	// No error
	return nil
}
