// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package from

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/compare"
	"zntr.io/harp/v2/pkg/tasks"
)

// OPLogTask implements secret-container creation from OpLog.
type OPLogTask struct {
	JSONReader   tasks.ReaderProvider
	OutputWriter tasks.WriterProvider
}

// Run the task.
func (t *OPLogTask) Run(ctx context.Context) error {
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

	// Convert input as a map
	var input compare.OpLog
	if err = json.NewDecoder(reader).Decode(&input); err != nil {
		return fmt.Errorf("unable to decode input JSON: %w", err)
	}

	// Build the container from json oplog
	b, err = bundle.FromOpLog(input)
	if err != nil {
		return fmt.Errorf("unable to create container from oplog: %w", err)
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
