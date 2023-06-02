// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package from

import (
	"context"
	"fmt"
	"io"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/hcl"
	"zntr.io/harp/v2/pkg/tasks"
)

// JSONMapTask implements secret-container creation from JSON Map.
type HCLTask struct {
	HCLReader    tasks.ReaderProvider
	OutputWriter tasks.WriterProvider
}

// Run the task.
func (t *HCLTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		writer io.Writer
		b      *bundlev1.Bundle
		err    error
	)

	// Create input reader
	reader, err = t.HCLReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to read input reader: %w", err)
	}

	// Parse input as HCL configuration object.
	cfg, err := hcl.Parse(reader, "input", "hcl")
	if err != nil {
		return fmt.Errorf("unable to parse input HCL: %w", err)
	}

	// Build the container from hcl dsl
	b, err = bundle.FromHCL(cfg)
	if err != nil {
		return fmt.Errorf("unable to create container from hcl: %w", err)
	}

	// Create output writer
	writer, err = t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output writer: %w", err)
	}

	// Dump bundle
	if err = bundle.ToContainerWriter(writer, b); err != nil {
		return fmt.Errorf("unable to produce exported bundle: %w", err)
	}

	// No error
	return nil
}
