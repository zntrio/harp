// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package yaml

import (
	"testing"
)

func TestYamlParser(t *testing.T) {
	parser := &Parser{}
	sample := `foo:
  test:
    bar: "toto"`

	var input interface{}
	if err := parser.Unmarshal([]byte(sample), &input); err != nil {
		t.Fatalf("parser should not have thrown an error: %v", err)
	}

	if input == nil {
		t.Fatalf("there should be information parsed but its nil")
	}

	inputMap := input.(map[string]interface{})
	item := inputMap["foo"]
	if len(item.(map[string]interface{})) <= 0 {
		t.Error("there should be at least one item defined in the parsed file, but none found")
	}
}
