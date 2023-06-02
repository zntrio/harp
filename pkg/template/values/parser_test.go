// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package values

import (
	"reflect"
	"testing"

	"zntr.io/harp/v2/pkg/template/values/hcl1"
	"zntr.io/harp/v2/pkg/template/values/hcl2"
	"zntr.io/harp/v2/pkg/template/values/hocon"
	"zntr.io/harp/v2/pkg/template/values/toml"
	"zntr.io/harp/v2/pkg/template/values/xml"
	"zntr.io/harp/v2/pkg/template/values/yaml"
)

func TestGetParser(t *testing.T) {
	testTable := []struct {
		name        string
		fileType    string
		expected    Parser
		expectError bool
	}{
		{
			name:        "Test getting HOCON parser",
			fileType:    "hocon",
			expected:    new(hocon.Parser),
			expectError: false,
		},
		{
			name:        "Test getting TOML parser",
			fileType:    "toml",
			expected:    new(toml.Parser),
			expectError: false,
		},
		{
			name:        "Test getting XML parser",
			fileType:    "xml",
			expected:    new(xml.Parser),
			expectError: false,
		},
		{
			name:        "Test getting Terraform parser from HCL1 input",
			fileType:    "hcl1",
			expected:    new(hcl1.Parser),
			expectError: false,
		},
		{
			name:        "Test getting Terraform parser from HCL2 input",
			fileType:    "tf",
			expected:    new(hcl2.Parser),
			expectError: false,
		},
		{
			name:        "Test getting Terraform parser from YAML input",
			fileType:    "yaml",
			expected:    new(yaml.Parser),
			expectError: false,
		},
		{
			name:        "Test getting Terraform parser from JSON input",
			fileType:    "json",
			expected:    new(yaml.Parser),
			expectError: false,
		},
		{
			name:        "Test getting invalid filetype",
			fileType:    "epicfailure",
			expected:    nil,
			expectError: true,
		},
	}

	for _, testUnit := range testTable {
		t.Run(testUnit.name, func(t *testing.T) {
			received, err := GetParser(testUnit.fileType)

			if !reflect.DeepEqual(received, testUnit.expected) {
				t.Errorf("expected: %T \n got this: %T", testUnit.expected, received)
			}
			if !testUnit.expectError && err != nil {
				t.Errorf("error here: %v", err)
			}
			if testUnit.expectError && err == nil {
				t.Error("error expected but not received")
			}
		})
	}
}
