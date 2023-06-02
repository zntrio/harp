// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hcl1

import (
	"fmt"

	"github.com/hashicorp/hcl"
)

// Parser is a parser for TF and HCL files.
type Parser struct{}

// Unmarshal unmarshals TF and HCL files.
func (s *Parser) Unmarshal(p []byte, v interface{}) error {
	if err := hcl.Unmarshal(p, v); err != nil {
		return fmt.Errorf("unmarshal hcl: %w", err)
	}

	return nil
}
