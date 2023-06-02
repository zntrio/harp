// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ruleset

import (
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/protojson"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/sdk/convert"
	"zntr.io/harp/v2/pkg/sdk/types"
)

// YAML a given reader in order to extract a BundlePatch sepcification.
func YAML(r io.Reader) (*bundlev1.RuleSet, error) {
	// Check arguments
	if types.IsNil(r) {
		return nil, fmt.Errorf("reader is nil")
	}

	// Drain the reader
	jsonReader, err := convert.YAMLtoJSON(r)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input as BundlePatch: %w", err)
	}

	// Drain reader
	jsonData, err := io.ReadAll(jsonReader)
	if err != nil {
		return nil, fmt.Errorf("unable to drain all json reader content: %w", err)
	}

	// Initialize empty definition object
	def := bundlev1.RuleSet{}
	def.Reset()

	// Deserialize JSON with JSONPB wrapper
	if err := protojson.Unmarshal(jsonData, &def); err != nil {
		return nil, fmt.Errorf("unable to decode spec as json: %w", err)
	}

	// Validate spec
	if err := Validate(&def); err != nil {
		return nil, fmt.Errorf("unable to validate descriptor: %w", err)
	}

	// No error
	return &def, nil
}
