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
	"zntr.io/harp/v2/pkg/tasks"
)

// BundleDumpTask implements secret-container creation from a Bundle Dump.
type BundleDumpTask struct {
	JSONReader   tasks.ReaderProvider
	OutputWriter tasks.WriterProvider
}

// Run the task.
func (t *BundleDumpTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		writer io.Writer
		b      *bundlev1.Bundle
		err    error
	)

	// Create input reader
	reader, err = t.JSONReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to read input reader: %w", err)
	}

	// Build the container from json
	b, err = bundle.FromDump(reader)
	if err != nil {
		return fmt.Errorf("unable to create container from dump: %w", err)
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
