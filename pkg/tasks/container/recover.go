// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"zntr.io/harp/v2/pkg/container/identity"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/tasks"
)

// RecoverTask implements secret container identity recovery task.
type RecoverTask struct {
	JSONReader   tasks.ReaderProvider
	OutputWriter tasks.WriterProvider
	Transformer  value.Transformer
	JSONOutput   bool
}

// Run the task.
func (t *RecoverTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.JSONReader) {
		return errors.New("unable to run task with a nil jsonReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if types.IsNil(t.Transformer) {
		return errors.New("unable to run task with a nil transformer")
	}

	// Create input reader
	reader, err := t.JSONReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to read input reader: %w", err)
	}

	// Extract from reader
	input, err := identity.FromReader(reader)
	if err != nil {
		return fmt.Errorf("unable to extract an identity from reader: %w", err)
	}

	// Try to decrypt the private key
	privateKey, err := input.Decrypt(ctx, t.Transformer)
	if err != nil {
		return fmt.Errorf("unable to decrypt private key: %w", err)
	}

	// Retrieve recovery key
	recoveryPrivateKey, err := privateKey.RecoveryKey()
	if err != nil {
		return fmt.Errorf("unable to retrieve recovery key from identity: %w", err)
	}

	// Get output writer
	outputWriter, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve output writer: %w", err)
	}

	// Display as json
	if t.JSONOutput {
		if errJSON := json.NewEncoder(outputWriter).Encode(map[string]interface{}{
			"container_key": recoveryPrivateKey,
		}); errJSON != nil {
			return fmt.Errorf("unable to display as json: %w", errJSON)
		}
	} else {
		// Display container key
		if _, err := fmt.Fprintf(outputWriter, "Container key : %s\n", recoveryPrivateKey); err != nil {
			return fmt.Errorf("unable to display result: %w", err)
		}
	}

	// No error
	return nil
}
