// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package yaml

import (
	"fmt"

	"sigs.k8s.io/yaml"
)

// Parser is an XML parser.
type Parser struct{}

// Unmarshal unmarshals YAML files.
func (p *Parser) Unmarshal(body []byte, v interface{}) error {
	if err := yaml.Unmarshal(body, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}
