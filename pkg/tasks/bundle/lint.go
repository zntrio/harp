// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"errors"
	"fmt"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/ruleset"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// LintTask implements bundle linting task.
type LintTask struct {
	ContainerReader tasks.ReaderProvider
	RuleSetReader   tasks.ReaderProvider
}

// Run the task.
func (t *LintTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.RuleSetReader) {
		return errors.New("unable to run task with a nil ruleSetReader provider")
	}

	// Create input reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize bundle reader: %w", err)
	}

	// Create input reader
	rsReader, err := t.RuleSetReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize ruleset reader: %w", err)
	}

	// Parse the input specification
	spec, err := ruleset.YAML(rsReader)
	if err != nil {
		return fmt.Errorf("unable to parse ruleset file: %w", err)
	}

	// Load bundle
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle content: %w", err)
	}

	if err := ruleset.Evaluate(ctx, b, spec); err != nil {
		return fmt.Errorf("unable to validate given bundle: %w", err)
	}

	// No error
	return nil
}
