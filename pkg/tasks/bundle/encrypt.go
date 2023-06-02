// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"errors"
	"fmt"
	"io"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/tasks"
)

// EncryptTask implements secret container encryption task.
type EncryptTask struct {
	ContainerReader   tasks.ReaderProvider
	OutputWriter      tasks.WriterProvider
	BundleTransformer value.Transformer
	TransformerMap    map[string]value.Transformer
	SkipUnresolved    bool
}

// Run the task.
func (t *EncryptTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		writer io.Writer
		b      *bundlev1.Bundle
		err    error
	)

	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}

	// Create input reader
	reader, err = t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle: %w", err)
	}

	// Read input bundle
	b, err = bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to read input as bundle: %w", err)
	}

	// Select appropriate encryption strategy.
	switch {
	case !types.IsNil(t.BundleTransformer):
		// Apply transformer to bundle
		if err = bundle.Lock(ctx, b, t.BundleTransformer); err != nil {
			return fmt.Errorf("unable to apply bundle transformation: %w", err)
		}
	case len(t.TransformerMap) > 0:
		// Apply annotation based encryption
		if err = bundle.PartialLock(ctx, b, t.TransformerMap, t.SkipUnresolved); err != nil {
			return fmt.Errorf("unable to apply annotation based transformation: %w", err)
		}
	default:
		return errors.New("invalid encryption strategy, can't determine if it's a full bundle or a selective annotation based encryption")
	}

	// Create output writer
	writer, err = t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output bundle: %w", err)
	}

	// Dump bundle
	if err = bundle.ToContainerWriter(writer, b); err != nil {
		return fmt.Errorf("unable to produce transformed bundle: %w", err)
	}

	// No error
	return nil
}
