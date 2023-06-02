// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"errors"
	"fmt"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/patch"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// PatchTask implements secret container patching task.
type PatchTask struct {
	PatchReader     tasks.ReaderProvider
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	Values          map[string]interface{}
	Options         []patch.OptionFunc
}

// Run the task.
func (t *PatchTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.PatchReader) {
		return errors.New("unable to run task with a nil patchReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}

	// Retrieve the container reader
	containerReader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve patch reader: %w", err)
	}

	// Load bundle
	b, err := bundle.FromContainerReader(containerReader)
	if err != nil {
		return fmt.Errorf("unable to load bundle content: %w", err)
	}

	// Retrieve the patch reader
	patchReader, err := t.PatchReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve patch reader: %w", err)
	}

	// Parse the input specification
	spec, err := patch.YAML(patchReader)
	if err != nil {
		return fmt.Errorf("unable to parse patch file: %w", err)
	}

	// Apply the patch speicification to generate an output bundle
	patchedBundle, err := patch.Apply(ctx, spec, b, t.Values, t.Options...)
	if err != nil {
		return fmt.Errorf("unable to generate output bundle from patch: %w", err)
	}

	// Retrieve the container reader
	outputWriter, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve output writer: %w", err)
	}

	// Dump all content
	if err = bundle.ToContainerWriter(outputWriter, patchedBundle); err != nil {
		return fmt.Errorf("unable to dump bundle content: %w", err)
	}

	// No error
	return nil
}
