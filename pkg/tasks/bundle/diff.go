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
	"zntr.io/harp/v2/pkg/bundle/compare"
	"zntr.io/harp/v2/pkg/sdk/convert"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// DiffTask implements secret container difference task.
type DiffTask struct {
	SourceReader      tasks.ReaderProvider
	DestinationReader tasks.ReaderProvider
	OutputWriter      tasks.WriterProvider
	GeneratePatch     bool
}

// Run the task.
func (t *DiffTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.SourceReader) {
		return errors.New("unable to run task with a nil sourceReader provider")
	}
	if types.IsNil(t.DestinationReader) {
		return errors.New("unable to run task with a nil destinationReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}

	// Create input reader
	readerSrc, err := t.SourceReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open source bundle: %w", err)
	}

	// Load source bundle
	bSrc, err := bundle.FromContainerReader(readerSrc)
	if err != nil {
		return fmt.Errorf("unable to load source bundle content: %w", err)
	}

	// Create input reader
	readerDst, err := t.DestinationReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open destination bundle: %w", err)
	}

	// Load destination bundle
	bDst, err := bundle.FromContainerReader(readerDst)
	if err != nil {
		return fmt.Errorf("unable to load destination bundle content: %w", err)
	}

	// Calculate diff
	report, err := compare.Diff(bSrc, bDst)
	if err != nil {
		return fmt.Errorf("unable to calculate bundle difference: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output writer: %w", err)
	}

	if !t.GeneratePatch {
		// Encode as JSON
		if err := json.NewEncoder(writer).Encode(report); err != nil {
			return fmt.Errorf("unable to marshal JSON OpLog: %w", err)
		}
	} else {
		// Convert optlog as a patch
		patch, err := compare.ToPatch(report)
		if err != nil {
			return fmt.Errorf("unable to convert oplog as a bundle patch: %w", err)
		}

		// Marshal as YAML
		out, err := convert.PBtoYAML(patch)
		if err != nil {
			return fmt.Errorf("unable to marshal patch as YAML: %w", err)
		}

		// Write output
		fmt.Fprintln(writer, string(out))
	}

	// No error
	return nil
}
