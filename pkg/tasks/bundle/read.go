// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// ReadTask implements secret container reading task.
type ReadTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	PackageName     string
	SecretKey       string
}

// Run the task.
func (t *ReadTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if t.PackageName == "" {
		return errors.New("unable to proceed with blank packageName")
	}

	// Create input reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle: %w", err)
	}

	// Load bundle
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle content: %w", err)
	}

	// Read a secret from bundle
	s, err := bundle.Read(b, t.PackageName)
	if err != nil {
		return fmt.Errorf("unable to read bundle content: %w", err)
	}

	// Prepare output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to get output writer: %w", err)
	}

	if t.SecretKey != "" {
		if v, ok := s[t.SecretKey]; ok {
			fmt.Fprintf(writer, "%s", v)
		} else {
			return fmt.Errorf("requested field does not exist %q: %w", t.SecretKey, err)
		}
	} else {
		// Dump the secret value
		if err := json.NewEncoder(writer).Encode(s); err != nil {
			return fmt.Errorf("unable to encode secret value as json: %w", err)
		}
	}

	// No error
	return nil
}
