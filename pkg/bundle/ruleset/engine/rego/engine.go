// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package rego

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/open-policy-agent/opa/rego"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/ruleset/engine"
)

const (
	maxPolicySize = 5 * 1024 * 1025 // 5MB
)

func New(ctx context.Context, r io.Reader) (engine.PackageLinter, error) {
	// Read all policy content
	policy, err := io.ReadAll(io.LimitReader(r, maxPolicySize))
	if err != nil {
		return nil, fmt.Errorf("unable to read the policy content: %w", err)
	}

	// Parse and prepare the policy
	query, err := rego.New(
		rego.Query("data.harp.compliant"),
		rego.Module("harp.rego", string(policy)),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare for eval: %w", err)
	}

	// Return engine
	return &ruleEngine{
		query: query,
	}, nil
}

// -----------------------------------------------------------------------------

type ruleEngine struct {
	query rego.PreparedEvalQuery
}

func (re *ruleEngine) EvaluatePackage(ctx context.Context, p *bundlev1.Package) error {
	// Check arguments
	if p == nil {
		return errors.New("unable to evaluate nil package")
	}

	// Evaluation with the given package
	results, err := re.query.Eval(ctx, rego.EvalInput(p))
	if err != nil {
		return fmt.Errorf("unable to evaluate the policy: %w", err)
	} else if len(results) == 0 {
		// Handle undefined result.
		return nil
	}

	for _, result := range results {
		for _, expression := range result.Expressions {
			// Extract result
			compliant, ok := expression.Value.(bool)
			if !ok {
				// Handle unexpected result type.
				return errors.New("the policy must return boolean")
			}

			// Check package compliance
			if !compliant {
				return engine.ErrRuleNotValid
			}
		}
	}

	// Package validated
	return nil
}
