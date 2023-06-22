// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/awnumar/memguard"
	"zntr.io/harp/v2/pkg/container"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// UnsealTask implements secret container unsealing task.
type UnsealTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	ContainerKey    *memguard.LockedBuffer
	PreSharedKey    *memguard.LockedBuffer
}

// Run the task.
func (t *UnsealTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if t.ContainerKey == nil {
		return errors.New("unable to run task with a nil container key")
	}

	// Create input reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle reader: %w", err)
	}

	// Load input container
	in, err := container.Load(reader)
	if err != nil {
		return fmt.Errorf("unable to read input container: %w", err)
	}

	// Seal options
	sopts := []container.Option{}

	// Process pre-shared key
	if t.PreSharedKey != nil {
		// Try to decode preshared key
		psk, errDecode := base64.RawURLEncoding.DecodeString(t.PreSharedKey.String())
		if errDecode != nil {
			return fmt.Errorf("unable to decode pre-shared key: %w", errDecode)
		}
		sopts = append(sopts, container.WithPreSharedKey(memguard.NewBufferFromBytes(psk)))
		t.PreSharedKey.Destroy()
	}

	// Unseal the container
	out, err := container.Unseal(in, t.ContainerKey, sopts...)
	if err != nil {
		return fmt.Errorf("unable to unseal bundle content: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output bundle: %w", err)
	}

	// Dump all content
	if err := container.Dump(writer, out); err != nil {
		return fmt.Errorf("unable to write unsealed container: %w", err)
	}

	// No error
	return nil
}
