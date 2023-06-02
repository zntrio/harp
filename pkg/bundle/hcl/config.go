// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hcl

// Config is the configuration structure for bundle DSL.
type Config struct {
	Annotations map[string]string `hcl:"annotations,optional"`
	Labels      map[string]string `hcl:"labels,optional"`
	Packages    []Package         `hcl:"package,block"`
}

type Package struct {
	Path        string            `hcl:"path,label"`
	Description string            `hcl:"description"`
	Annotations map[string]string `hcl:"annotations,optional"`
	Labels      map[string]string `hcl:"labels,optional"`
	Secrets     map[string]string `hcl:"secrets"`
}
