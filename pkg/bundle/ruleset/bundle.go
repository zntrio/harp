// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ruleset

import (
	"encoding/base64"
	"errors"
	"fmt"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
)

// FromBundle crawls secret structure to generate a linter ruleset.
func FromBundle(b *bundlev1.Bundle) (*bundlev1.RuleSet, error) {
	// Check arguments
	if b == nil {
		return nil, errors.New("unable to process nil bundle")
	}
	if len(b.Packages) == 0 {
		return nil, errors.New("unable to generate rule from an empty bundle")
	}

	// Retrieve MTR
	root, _, err := bundle.Tree(b)
	if err != nil {
		return nil, fmt.Errorf("unable to compute bundle identifier: %w", err)
	}

	// Encode MTR as Base64
	b64Root := base64.RawURLEncoding.EncodeToString(root.Root())

	// Create ruleset
	rs := &bundlev1.RuleSet{
		ApiVersion: "harp.zntr.io/v2",
		Kind:       "RuleSet",
		Meta: &bundlev1.RuleSetMeta{
			Name:        b64Root,
			Description: "Generated from bundle content",
		},
		Spec: &bundlev1.RuleSetSpec{
			Rules: []*bundlev1.Rule{},
		},
	}

	// Iterate over bundle package
	ruleIdx := 1
	for _, p := range b.Packages {
		if p == nil || p.Secrets == nil || len(p.Secrets.Data) == 0 {
			// Skip invalid package
			continue
		}

		// Prepare a rule
		r := &bundlev1.Rule{
			Name:        fmt.Sprintf("LINT-%s-%d", b64Root[:6], ruleIdx),
			Path:        p.Name,
			Constraints: []string{},
		}

		// Process each secret
		for _, s := range p.Secrets.Data {
			r.Constraints = append(r.Constraints, fmt.Sprintf(`p.has_secret(%q)`, s.Key))
		}

		// Add the rules
		rs.Spec.Rules = append(rs.Spec.Rules, r)
		ruleIdx++
	}

	// No error
	return rs, nil
}
