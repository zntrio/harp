// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package to

import (
	"context"
	"fmt"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/ruleset"
	"zntr.io/harp/v2/pkg/sdk/convert"
	"zntr.io/harp/v2/pkg/tasks"
)

// RuleSetTask implements RuleSet generation from a bundle.
type RuleSetTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
}

// Run the task.
func (t *RuleSetTask) Run(ctx context.Context) error {
	// Create input reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize bundle reader: %w", err)
	}

	// Load bundle
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle content: %w", err)
	}

	// Generate ruleset
	rs, err := ruleset.FromBundle(b)
	if err != nil {
		return fmt.Errorf("unable to generate RuleSet from given bundle: %w", err)
	}

	// Marshal as YAML
	out, err := convert.PBtoYAML(rs)
	if err != nil {
		return fmt.Errorf("unable to marshal patch as YAML: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize output writer: %w", err)
	}

	// Write output
	fmt.Fprintln(writer, string(out))

	// No error
	return nil
}
