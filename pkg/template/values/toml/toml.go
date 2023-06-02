// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package toml

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

// Parser is a TOML parser.
type Parser struct{}

// Unmarshal unmarshals TOML files.
func (tp *Parser) Unmarshal(p []byte, v interface{}) error {
	if err := toml.Unmarshal(p, v); err != nil {
		return fmt.Errorf("unmarshal toml: %w", err)
	}

	return nil
}
