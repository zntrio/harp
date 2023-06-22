// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ruleset

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"

	"golang.org/x/crypto/blake2b"
)

// Validate bundle patch.
func Validate(spec *bundlev1.RuleSet) error {
	// Check if spec is nil
	if spec == nil {
		return fmt.Errorf("unable to validate bundle patch: patch is nil")
	}

	if spec.ApiVersion != "harp.zntr.io/v2" {
		return fmt.Errorf("apiVersion should be 'harp.zntr.io/v2'")
	}

	if spec.Kind != "RuleSet" {
		return fmt.Errorf("kind should be 'RuleSet'")
	}

	if spec.Meta == nil {
		return fmt.Errorf("meta should be 'nil'")
	}

	if spec.Spec == nil {
		return fmt.Errorf("spec should be 'nil'")
	}

	// No error
	return nil
}

// Checksum calculates the specification checksum.
func Checksum(spec *bundlev1.RuleSet) (string, error) {
	// Validate bundle template
	if err := Validate(spec); err != nil {
		return "", fmt.Errorf("unable to validate spec: %w", err)
	}

	// Encode spec as protobuf
	payload, err := proto.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("unable to encode bundle patch: %w", err)
	}

	// Calculate checksum
	checksum := blake2b.Sum256(payload)

	// No error
	return base64.RawURLEncoding.EncodeToString(checksum[:]), nil
}
