// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hcl

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// ParseFile parses the given file for a configuration. The syntax of the
// file is determined based on the filename extension: "hcl" for HCL,
// "json" for JSON, other is an error.
func ParseFile(filename string) (*Config, error) {
	var config Config
	return &config, hclsimple.DecodeFile(filename, nil, &config)
}

// Parse parses the configuration from the given reader. The reader will be
// read to completion (EOF) before returning so ensure that the reader
// does not block forever.
//
// format is either "hcl" or "json".
func Parse(r io.Reader, filename, format string) (*Config, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to drain input reader: %w", err)
	}

	var config Config
	return &config, hclsimple.Decode("config.hcl", src, nil, &config)
}
