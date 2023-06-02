// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// PrefixerTask implements secret container prefix management task.
type PrefixerTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	Prefix          string
	Remove          bool
}

// Run the task.
func (t *PrefixerTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if t.Prefix == "" {
		return errors.New("unable to proceed with blank prefix")
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

	// Iterate over all packages
	for _, p := range b.Packages {
		if t.Remove {
			p.Name = strings.TrimPrefix(path.Clean(strings.TrimPrefix(p.Name, t.Prefix)), "/")
		} else {
			p.Name = path.Clean(path.Join(t.Prefix, p.Name))
		}
	}

	// Retrieve the container reader
	outputWriter, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve output writer: %w", err)
	}

	// Dump all content
	if err = bundle.ToContainerWriter(outputWriter, b); err != nil {
		return fmt.Errorf("unable to dump bundle content: %w", err)
	}

	// No error
	return nil
}
