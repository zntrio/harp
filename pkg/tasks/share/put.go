// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package share

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"zntr.io/harp/v2/pkg/tasks"
	"zntr.io/harp/v2/pkg/vault"
)

// PutTask implements secret sharing via Vault Cubbyhole.
type PutTask struct {
	InputReader    tasks.ReaderProvider
	OutputWriter   tasks.WriterProvider
	BackendPrefix  string
	TTL            time.Duration
	VaultNamespace string
	JSONOutput     bool
}

// Run the task.
func (t *PutTask) Run(ctx context.Context) error {
	// Create input reader
	reader, err := t.InputReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to read input reader: %w", err)
	}

	// Initialize vault connection
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("unable to initialize Vault connection: %w", err)
	}

	// If a namespace is specified
	if t.VaultNamespace != "" {
		client.SetNamespace(t.VaultNamespace)
	}

	// Set expiration
	client.SetWrappingLookupFunc(func(operation, path string) string {
		return t.TTL.String()
	})

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

	// Get output writer
	outputWriter, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve output writer: %w", err)
	}

	// Retrieve secret
	token, err := s.Put(ctx, reader)
	if err != nil {
		return fmt.Errorf("unable to retrieve secret: %w", err)
	}

	// Display as json
	if t.JSONOutput {
		if err := json.NewEncoder(outputWriter).Encode(map[string]interface{}{
			"token":      token,
			"expires_in": t.TTL.Seconds(),
		}); err != nil {
			return fmt.Errorf("unable to display as json: %w", err)
		}
	} else {
		// Display container key
		fmt.Fprintf(outputWriter, "Token : %s (Expires in %d seconds)\n", token, int64(t.TTL.Seconds()))
	}

	// No error
	return nil
}
