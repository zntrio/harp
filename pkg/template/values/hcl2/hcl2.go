// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hcl2

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Parser is a HCL2 parser.
type Parser struct{}

// Unmarshal HCL2.0 scripts.
func (h *Parser) Unmarshal(p []byte, v interface{}) error {
	file, diags := hclsyntax.ParseConfig(p, "", hcl.Pos{Byte: 0, Line: 1, Column: 1})

	if diags.HasErrors() {
		var details []error
		details = append(details, diags.Errs()...)
		return fmt.Errorf("parse hcl2 config: \n %s", details)
	}

	content, err := convertFile(file)
	if err != nil {
		return fmt.Errorf("convert hcl2 to json: %w", err)
	}

	j, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal hcl2 to json: %w", err)
	}

	if err := json.Unmarshal(j, v); err != nil {
		return fmt.Errorf("unmarshal hcl2 json: %w", err)
	}

	return nil
}
