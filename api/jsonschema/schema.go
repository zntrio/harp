// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jsonschema

import (
	_ "embed"
)

//go:embed harp.bundle.v1/Bundle.json
var bundleV1BundleSchemaDefinition []byte

// BundleV1BundleSchema returns the `harp.bundle.v1.Bundle` jsonschema content.
func BundleV1BundleSchema() []byte {
	return bundleV1BundleSchemaDefinition
}

//go:embed harp.bundle.v1/Patch.json
var bundleV1PatchSchemaDefinition []byte

// BundleV1PatchSchema returns the `harp.bundle.v1.Patch` jsonschema content.
func BundleV1PatchSchema() []byte {
	return bundleV1PatchSchemaDefinition
}

//go:embed harp.bundle.v1/RuleSet.json
var bundleV1RuleSetSchemaDefinition []byte

// BundleV1RuleSetSchema returns the `harp.bundle.v1.RuleSet` jsonschema content.
func BundleV1RuleSetSchema() []byte {
	return bundleV1RuleSetSchemaDefinition
}

//go:embed harp.bundle.v1/Template.json
var bundleV1TemplateSchemaDefinition []byte

// BundleV1TemplateSchema returns the `harp.bundle.v1.Template` jsonschema content.
func BundleV1TemplateSchema() []byte {
	return bundleV1TemplateSchemaDefinition
}
