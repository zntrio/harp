// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package template

import (
	"errors"
	"fmt"
	"io"

	"github.com/xeipuuv/gojsonschema"
	"zntr.io/harp/v2/api/jsonschema"
	"zntr.io/harp/v2/pkg/sdk/convert"
	"zntr.io/harp/v2/pkg/sdk/types"
)

// JSONSchema returns the used json schema for validation.
func JSONSchema() []byte {
	return jsonschema.BundleV1TemplateSchema()
}

// Lint to input reader content with Bundle jsonschema.
func Lint(r io.Reader) ([]gojsonschema.ResultError, error) {
	// Check arguments
	if types.IsNil(r) {
		return nil, fmt.Errorf("reader is nil")
	}

	// Drain the reader
	jsonReader, err := convert.YAMLtoJSON(r)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input as YAML: %w", err)
	}

	// Drain reader
	jsonData, err := io.ReadAll(jsonReader)
	if err != nil {
		return nil, fmt.Errorf("unable to drain all json reader content: %w", err)
	}

	// Prepare loaders
	schemaLoader := gojsonschema.NewBytesLoader(jsonschema.BundleV1TemplateSchema())
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("template validation failed %w", err)
	}
	if !result.Valid() {
		return result.Errors(), errors.New("template not valid")
	}

	// No error
	return nil, nil
}
