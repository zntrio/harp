// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package from

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/value/flatmap"
	"zntr.io/harp/v2/pkg/tasks"
)

// ObjectTask implements secret-container creation from a YAML/JSON structure.
type ObjectTask struct {
	ObjectReader tasks.ReaderProvider
	OutputWriter tasks.WriterProvider
	JSON         bool
	YAML         bool
}

// Run the task.
func (t *ObjectTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		writer io.Writer
		b      *bundlev1.Bundle
		err    error
	)

	// Create input reader
	reader, err = t.ObjectReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to read input reader: %w", err)
	}

	// Decode as YAML any object
	var source map[string]interface{}

	switch {
	case t.YAML:
		if errYaml := yaml.NewDecoder(reader).Decode(&source); errYaml != nil {
			return fmt.Errorf("unable to decode source as YAML: %w", err)
		}
	case t.JSON:
		if errYaml := json.NewDecoder(reader).Decode(&source); errYaml != nil {
			return fmt.Errorf("unable to decode source as JSON: %w", err)
		}
	default:
		return errors.New("json or yaml must be selected")
	}

	// Flatten the struct
	input := flatmap.Flatten(source)

	// Build the container from json
	b, err = bundle.FromMap(input)
	if err != nil {
		return fmt.Errorf("unable to create container from map: %w", err)
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
