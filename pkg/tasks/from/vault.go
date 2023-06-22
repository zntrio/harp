// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package from

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"zntr.io/harp/v2/pkg/bundle"
	bundlevault "zntr.io/harp/v2/pkg/bundle/vault"
	"zntr.io/harp/v2/pkg/tasks"
	"zntr.io/harp/v2/pkg/vault"
)

// VaultTask implements secret-container building from Vault K/V.
type VaultTask struct {
	OutputWriter    tasks.WriterProvider
	SecretPaths     []string
	VaultNamespace  string
	AsVaultMetadata bool
	WithMetadata    bool
	MaxWorkerCount  int64
	ContinueOnError bool
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

	// Call exporter
	b, err := bundlevault.Pull(ctx, client, t.SecretPaths,
		bundlevault.WithVaultMetadata(t.AsVaultMetadata),
		bundlevault.WithSecretMetadata(t.WithMetadata),
		bundlevault.WithMaxWorkerCount(t.MaxWorkerCount),
		bundlevault.WithContinueOnError(t.ContinueOnError),
	)
	if err != nil {
		return fmt.Errorf("error occurs during vault export: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output bundle: %w", err)
	}

	// Dump bundle
	if err = bundle.ToContainerWriter(writer, b); err != nil {
		return fmt.Errorf("unable to produce exported bundle: %w", err)
	}

	// No error
	return nil
}
