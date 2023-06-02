// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package values

import (
	"fmt"

	"zntr.io/harp/v2/pkg/template/values/hcl1"
	"zntr.io/harp/v2/pkg/template/values/hcl2"
	"zntr.io/harp/v2/pkg/template/values/hocon"
	"zntr.io/harp/v2/pkg/template/values/toml"
	"zntr.io/harp/v2/pkg/template/values/xml"
	"zntr.io/harp/v2/pkg/template/values/yaml"
)

// Parser is the interface implemented by objects that can unmarshal
// bytes into a golang interface.
type Parser interface {
	Unmarshal(p []byte, v interface{}) error
}

// GetParser gets a file parser based on the file type and input.
func GetParser(fileType string) (Parser, error) {
	switch fileType {
	case "toml":
		return &toml.Parser{}, nil
	case "hocon":
		return &hocon.Parser{}, nil
	case "xml":
		return &xml.Parser{}, nil
	case "json", "yaml", "yml":
		return &yaml.Parser{}, nil
	case "hcl", "tf", "hcl2", "tfvars":
		return &hcl2.Parser{}, nil
	case "hcl1":
		return &hcl1.Parser{}, nil
	default:
		return nil, fmt.Errorf("unknown filetype given: %v", fileType)
	}
}
