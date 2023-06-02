// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package template

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/template/visitor"
	"zntr.io/harp/v2/pkg/sdk/types"

	"golang.org/x/crypto/blake2b"
)

// Validate bundle template.
func Validate(spec *bundlev1.Template) error {
	// Check if spec is nil
	if spec == nil {
		return fmt.Errorf("unable to validate bundle template: template is nil")
	}

	if spec.ApiVersion != "harp.zntr.io/v2" {
		return fmt.Errorf("apiVersion should be 'BundleTemplate'")
	}

	if spec.Kind != "BundleTemplate" {
		return fmt.Errorf("kind should be 'BundleTemplate'")
	}

	if spec.Meta == nil {
		return fmt.Errorf("meta should not be 'nil'")
	}

	if spec.Spec == nil {
		return fmt.Errorf("spec should not be 'nil'")
	}

	// No error
	return nil
}

// Checksum calculates the bundle template checksum.
func Checksum(spec *bundlev1.Template) (string, error) {
	// Check if spec is nil
	if spec == nil {
		return "", fmt.Errorf("unable to compute template checksum: template is nil")
	}

	// Validate bundle template
	if err := Validate(spec); err != nil {
		return "", fmt.Errorf("unable to validate spec: %w", err)
	}

	// Encode spec as protobuf
	payload, err := proto.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("unable to encode bundle template: %w", err)
	}

	// Calculate checksum
	checksum := blake2b.Sum256(payload)

	// No error
	return base64.RawURLEncoding.EncodeToString(checksum[:]), nil
}

// Execute a template to generate a final secret bundle.
func Execute(spec *bundlev1.Template, v visitor.TemplateVisitor) error {
	// Check if spec is nil
	if spec == nil {
		return fmt.Errorf("unable to execute bundle template: template is nil")
	}
	if types.IsNil(v) {
		return fmt.Errorf("unable to execute bundle template: visitor is nil")
	}

	// Validate bundle template
	if err := Validate(spec); err != nil {
		return fmt.Errorf("unable to validate spec: %w", err)
	}

	// Walk all namespaces
	v.Visit(spec)

	// Check error
	return v.Error()
}
