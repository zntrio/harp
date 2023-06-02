// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package xml

import (
	"bytes"
	"encoding/json"
	"fmt"

	x "github.com/basgys/goxml2json"
)

// Parser is an XML parser.
type Parser struct{}

// Unmarshal unmarshals XML files.
func (xml *Parser) Unmarshal(p []byte, v interface{}) error {
	res, err := x.Convert(bytes.NewReader(p))
	if err != nil {
		return fmt.Errorf("unmarshal xml: %w", err)
	}

	if err := json.Unmarshal(res.Bytes(), v); err != nil {
		return fmt.Errorf("convert xml to json: %w", err)
	}

	return nil
}
