// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package template

import (
	"context"
	"encoding/json"
	"fmt"

	"zntr.io/harp/v2/pkg/tasks"
	tplcmdutil "zntr.io/harp/v2/pkg/template/cmdutil"
)

// ValueTask implements value object generation task.
type ValueTask struct {
	OutputWriter tasks.WriterProvider
	ValueFiles   []string
	Values       []string
	StringValues []string
	FileValues   []string
}

// Run the task.
func (t *ValueTask) Run(ctx context.Context) error {
	// Load values
	valueOpts := tplcmdutil.ValueOptions{
		ValueFiles:   t.ValueFiles,
		Values:       t.Values,
		StringValues: t.StringValues,
		FileValues:   t.FileValues,
	}
	values, err := valueOpts.MergeValues()
	if err != nil {
		return fmt.Errorf("unable to process values: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to create output writer: %w", err)
	}

	// Write rendered content
	if err := json.NewEncoder(writer).Encode(values); err != nil {
		return fmt.Errorf("unable to dump values as JSON: %w", err)
	}

	return nil
}
