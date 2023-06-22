// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/jmespath/go-jmespath"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// DumpTask implements secret-container dumping task.
type DumpTask struct {
	ContainerReader tasks.ReaderProvider
	OutputWriter    tasks.WriterProvider
	PathOnly        bool
	DataOnly        bool
	MetadataOnly    bool
	JMESPathFilter  string
	IgnoreTemplate  bool
}

// Run the task.
func (t *DumpTask) Run(ctx context.Context) error {
	var (
		reader io.Reader
		err    error
	)

	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerRedaer provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}

	// Create input reader
	reader, err = t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle: %w", err)
	}

	// Load bundle
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle content: %w", err)
	}

	// Create output writer
	writer, err := t.OutputWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to open writer: %w", err)
	}

	// Clean template if requested.
	if t.IgnoreTemplate {
		b.Template = nil
	}

	switch {
	case t.DataOnly:
		return t.dumpData(writer, b)
	case t.MetadataOnly:
		return t.dumpMetadata(writer, b)
	case t.PathOnly:
		return t.dumpPath(writer, b)
	case t.JMESPathFilter != "":
		return t.dumpFilter(writer, b)
	default:
		// Dump full structure.
		if err := bundle.AsProtoJSON(writer, b); err != nil {
			return fmt.Errorf("unable to generate JSON: %w", err)
		}
	}

	// No error
	return nil
}

func (t *DumpTask) dumpData(writer io.Writer, b *bundlev1.Bundle) error {
	// Check arguments
	if types.IsNil(writer) {
		return fmt.Errorf("unable to process nil writer")
	}
	if b == nil {
		return fmt.Errorf("unable to process nil bundle")
	}

	// Convert bundle as a map
	bMap, err := bundle.AsMap(b)
	if err != nil {
		return fmt.Errorf("unable to convert bundle content: %w", err)
	}

	// Encode as JSON
	if err := json.NewEncoder(writer).Encode(bMap); err != nil {
		return fmt.Errorf("unable to marshal JSON bundle content: %w", err)
	}

	return nil
}

func (t *DumpTask) dumpMetadata(writer io.Writer, b *bundlev1.Bundle) error {
	// Check arguments
	if types.IsNil(writer) {
		return fmt.Errorf("unable to process nil writer")
	}
	if b == nil {
		return fmt.Errorf("unable to process nil bundle")
	}

	// Export metadata as map
	metaMap, err := bundle.AsMetadataMap(b)
	if err != nil {
		return fmt.Errorf("unable to convert bundle content: %w", err)
	}

	// Encode as JSON
	if err := json.NewEncoder(writer).Encode(metaMap); err != nil {
		return fmt.Errorf("unable to marshal JSON bundle metadata: %w", err)
	}

	return nil
}

func (t *DumpTask) dumpPath(writer io.Writer, b *bundlev1.Bundle) error {
	// Check arguments
	if types.IsNil(writer) {
		return fmt.Errorf("unable to process nil writer")
	}
	if b == nil {
		return fmt.Errorf("unable to process nil bundle")
	}

	// Extract bundle paths
	paths, err := bundle.Paths(b)
	if err != nil {
		return fmt.Errorf("unable to extract bundle paths: %w", err)
	}

	// Print a xargs compatible list
	for _, p := range paths {
		_, err = fmt.Fprintf(writer, "%s\n", p)
		if err != nil {
			return fmt.Errorf("unable to write package path %q to stdout: %w", p, err)
		}
	}

	return nil
}

func (t *DumpTask) dumpFilter(writer io.Writer, b *bundlev1.Bundle) error {
	// Check arguments
	if types.IsNil(writer) {
		return fmt.Errorf("unable to process nil writer")
	}
	if b == nil {
		return fmt.Errorf("unable to process nil bundle")
	}

	// Filter bundle with JMESPath expression
	res, err := jmespath.Search(t.JMESPathFilter, b)
	if err != nil {
		return fmt.Errorf("unable to process JMESPath filter %q: %w", t.JMESPathFilter, err)
	}

	// Encode response
	out, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("unable to encode JMESPath filter result: %w", err)
	}

	// Write to writer
	if _, err := fmt.Fprintf(writer, "%s", string(out)); err != nil {
		return fmt.Errorf("unable to write JSON to stdout: %w", err)
	}

	return nil
}
