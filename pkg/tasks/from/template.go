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
	"zntr.io/harp/v2/pkg/bundle/template"
	"zntr.io/harp/v2/pkg/bundle/template/visitor/secretbuilder"
	"zntr.io/harp/v2/pkg/tasks"
	"zntr.io/harp/v2/pkg/template/engine"
)

// BundleTemplateTask implements secret-container generation from BundleTemplate
// manifest.
type BundleTemplateTask struct {
	TemplateReader  tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	TemplateContext engine.Context
}

// Run the task.
func (t *BundleTemplateTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		writer io.Writer
		spec   *bundlev1.Template
		err    error
	)

	// Create input reader
	reader, err = t.TemplateReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle template: %w", err)
	}

	// Parse the input specification
	spec, err = template.YAML(reader)
	if err != nil {
		return fmt.Errorf("unable to parse template: %w", err)
	}

	// Initialize output
	b := &bundlev1.Bundle{
		Template: spec,
	}

	// Initialize a bundle creator
	v := secretbuilder.New(b, t.TemplateContext)

	// Execute the template to generate an output bundle
	if err = template.Execute(spec, v); err != nil {
		return fmt.Errorf("unable to generate output bundle from template: %w", err)
	}

	// Create output writer
	writer, err = t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open output bundle: %w", err)
	}

	// Dump all content
	if err = bundle.ToContainerWriter(writer, b); err != nil {
		return fmt.Errorf("unable to dump bundle content: %w", err)
	}

	// No error
	return nil
}
