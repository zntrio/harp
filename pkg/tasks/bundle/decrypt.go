// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bundle

import (
	"context"
	"errors"
	"fmt"
	"io"

	bundlev1 "github.com/zntrio/harp/v2/api/gen/go/harp/bundle/v1"
	"github.com/zntrio/harp/v2/pkg/bundle"
	"github.com/zntrio/harp/v2/pkg/sdk/types"
	"github.com/zntrio/harp/v2/pkg/sdk/value"
	"github.com/zntrio/harp/v2/pkg/tasks"
)

// DecryptTask implements secret container decryption task.
type DecryptTask struct {
	ContainerReader    tasks.ReaderProvider
	OutputWriter       tasks.WriterProvider
	Transformers       []value.Transformer
	SkipNotDecryptable bool
}

// Run the task.
func (t *DecryptTask) Run(ctx context.Context) error {
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
	if len(t.Transformers) == 0 {
		return errors.New("unable to run task with an empty transformer list")
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

	// Apply transformer to bundle
	if err = bundle.UnLock(ctx, b, t.Transformers, t.SkipNotDecryptable); err != nil {
		return fmt.Errorf("unable to apply bundle transformation: %w", err)
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
