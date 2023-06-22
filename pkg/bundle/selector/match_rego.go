// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// MatchRego returns a Rego package matcher specification.
func MatchRego(ctx context.Context, policy string) (Specification, error) {
	// Prepare query filter
	query, err := rego.New(
		rego.Query("data.harp.matched"),
		rego.Module("harp.rego", policy),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare for eval: %w", err)
	}

	// Wrap as a builder
	return &regoMatcher{
		ctx:   ctx,
		query: query,
	}, nil
}

type regoMatcher struct {
	ctx   context.Context
	query rego.PreparedEvalQuery
}

// IsSatisfiedBy returns specification satisfaction status.
func (s *regoMatcher) IsSatisfiedBy(object interface{}) bool {
	// If object is a package
	if p, ok := object.(*bundlev1.Package); ok {
		// Evaluate filter compliance
		matched, err := s.regoEvaluate(s.ctx, s.query, p)
		if err != nil {
			log.For(s.ctx).Debug("rego evaluation failed", zap.Error(err))
			return false
		}

		return matched
	}

	return false
}

// -----------------------------------------------------------------------------

func (s *regoMatcher) regoEvaluate(ctx context.Context, query rego.PreparedEvalQuery, input interface{}) (bool, error) {
	// Evaluate the package with the policy
	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("unable to evaluate the policy: %w", err)
	} else if len(results) == 0 {
		// Handle undefined result.
		return false, nil
	}

	// Extract decision
	keep, ok := results[0].Expressions[0].Value.(bool)
	if !ok {
		// Handle unexpected result type.
		return false, errors.New("the policy must return boolean")
	}

	// No error
	return keep, nil
}
