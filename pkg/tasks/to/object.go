// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package to

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/value/flatmap"
	"zntr.io/harp/v2/pkg/tasks"
)

// ObjectTask implements secret-container publication process to json/yaml content.
type ObjectTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	Expand          bool
	TOML            bool
	YAML            bool
}

// Run the task.
func (t *ObjectTask) Run(ctx context.Context) error {
	// Create the reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle reader: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open writer: %w", err)
	}

	// Extract bundle from container
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle: %w", err)
	}

	// Convert as map
	bundleMap, err := bundle.AsMap(b)
	if err != nil {
		return fmt.Errorf("unable to transform the bundle as a map: %w", err)
	}

	var toEncode interface{}

	// Expand if required
	if t.Expand {
		toEncode = flatmap.Expand(bundleMap, "")
	} else {
		toEncode = bundleMap
	}

	// Select strategy
	switch {
	case t.YAML:
		// Encode as YAML
		if err := yaml.NewEncoder(writer).Encode(toEncode); err != nil {
			return fmt.Errorf("unable to marshal YAML bundle content: %w", err)
		}
	case t.TOML:
		// Encode as TOML
		if err := toml.NewEncoder(writer).Encode(toEncode); err != nil {
			return fmt.Errorf("unable to marshal TOML bundle content: %w", err)
		}
	default:
		// Encode as JSON
		if err := json.NewEncoder(writer).Encode(toEncode); err != nil {
			return fmt.Errorf("unable to marshal JSON bundle content: %w", err)
		}
	}

	// No error
	return nil
}
