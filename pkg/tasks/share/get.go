// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package share

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"

	"zntr.io/harp/v2/pkg/tasks"
	"zntr.io/harp/v2/pkg/vault"
)

// GetTask implements secret sharing via Vault Cubbyhole.
type GetTask struct {
	OutputWriter   tasks.WriterProvider
	BackendPrefix  string
	VaultNamespace string
	Token          string
}

// Run the task.
func (t *GetTask) Run(ctx context.Context) error {
	// Initialize vault connection
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("unable to initialize Vault connection: %w", err)
	}

	// If a namespace is specified
	if t.VaultNamespace != "" {
		client.SetNamespace(t.VaultNamespace)
	}

	// Verify vault connection
	if _, errAuth := vault.CheckAuthentication(ctx, client); errAuth != nil {
		return fmt.Errorf("vault connection verification failed: %w", errAuth)
	}

	// Create cubbyhole service
	sf, errFactory := vault.FromVaultClient(client)
	if err != nil {
		return fmt.Errorf("unable to initialize service factory: %w", errFactory)
	}
	s, errService := sf.Cubbyhole(t.BackendPrefix)
	if errService != nil {
		return fmt.Errorf("unable to initialize service factory: %w", errFactory)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output writer: %w", err)
	}

	// Retrieve secret
	if err := s.Get(ctx, t.Token, writer); err != nil {
		return fmt.Errorf("unable to retrieve secret: %w", err)
	}

	// No error
	return nil
}
