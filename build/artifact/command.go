// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package artifact

import (
	"github.com/iancoleman/strcase"
)

// Command specification.
type Command struct {
	Name        string
	Description string
	Package     string
	Module      string
	UseBoring   bool
}

// Kebab returns the kebab-case artifact name.
func (c Command) Kebab() string {
	return strcase.ToKebab(c.Name)
}

// Camel returns the CamelCase artifact name.
func (c Command) Camel() string {
	return strcase.ToCamel(c.Name)
}

// HasModule returns trus if artifact as a dedicated module.
func (c Command) HasModule() bool {
	return c.Module != ""
}
