// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package to

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"zntr.io/harp/v2/pkg/bundle"
	bundlevault "zntr.io/harp/v2/pkg/bundle/vault"
	"zntr.io/harp/v2/pkg/tasks"
	"zntr.io/harp/v2/pkg/vault"
)

// VaultTask implements secret-container publication process to Vault.
type VaultTask struct {
	ContainerReader tasks.ReaderProvider
	BackendPrefix   string
	PushMetadata    bool
	AsVaultMetadata bool
	VaultNamespace  string
	MaxWorkerCount  int64
}

// Run the task.
func (t *VaultTask) Run(ctx context.Context) error {
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

	// Create the reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle reader: %w", err)
	}

	// Extract bundle from container
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle: %w", err)
	}

	// Process push operation
	if err := bundlevault.Push(ctx, b, client,
		bundlevault.WithPrefix(t.BackendPrefix),
		bundlevault.WithSecretMetadata(t.PushMetadata),
		bundlevault.WithVaultMetadata(t.AsVaultMetadata),
		bundlevault.WithMaxWorkerCount(t.MaxWorkerCount),
	); err != nil {
		return fmt.Errorf("error occurs during vault export (prefix: %q): %w", t.BackendPrefix, err)
	}

	// No error
	return nil
}
