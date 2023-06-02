// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cel

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	celext "github.com/google/cel-go/ext"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/ruleset/engine"
	"zntr.io/harp/v2/pkg/bundle/ruleset/engine/cel/ext"
)

// -----------------------------------------------------------------------------

// New returns a Google CEL based linter engine.
func New(expressions []string) (engine.PackageLinter, error) {
	// Prepare CEL Environment
	env, err := cel.NewEnv(
		cel.Types(&bundlev1.Bundle{}, &bundlev1.Package{}, &bundlev1.SecretChain{}, &bundlev1.KV{}),
		ext.Packages(),
		ext.Secrets(),
		celext.Strings(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare CEL engine environment: %w", err)
	}

	// Assemble the complete ruleset
	ruleset := make([]cel.Program, 0, len(expressions))
	for _, exp := range expressions {
		// Parse expression
		parsed, issues := env.Parse(exp)
		if issues != nil && issues.Err() != nil {
			return nil, fmt.Errorf("unable to parse %q, go error: %w", exp, issues.Err())
		}

		// Extract AST
		ast, cerr := env.Check(parsed)
		if cerr != nil && cerr.Err() != nil {
			return nil, fmt.Errorf("invalid CEL expression: %w", cerr.Err())
		}

		// request matching is a boolean operation, so we don't really know
		// what to do if the expression returns a non-boolean type
		if ast.OutputType() != cel.BoolType {
			return nil, fmt.Errorf("CEL rule engine expects return type of bool, not %s", ast.OutputType())
		}

		// Compile the program
		p, err := env.Program(ast)
		if err != nil {
			return nil, fmt.Errorf("error while creating CEL program: %w", err)
		}

		// Add to context
		ruleset = append(ruleset, p)
	}

	// Return rule engine
	return &ruleEngine{
		cel:     env,
		ruleset: ruleset,
	}, nil
}

// -----------------------------------------------------------------------------

type ruleEngine struct {
	cel     *cel.Env
	ruleset []cel.Program
}

func (re *ruleEngine) EvaluatePackage(ctx context.Context, p *bundlev1.Package) error {
	// Check arguments
	if p == nil {
		return errors.New("unable to evaluate nil package")
	}

	// Apply evaluation (implicit AND between rules)
	for _, exp := range re.ruleset {
		// Evaluate using the bundle context
		out, _, err := exp.Eval(map[string]interface{}{
			"p": p,
		})
		if err != nil {
			return fmt.Errorf("an error occurred during the rule evaluation: %w", err)
		}

		// Boolean rule returned false
		if out.Value() == false {
			return engine.ErrRuleNotValid
		}
	}

	// No error
	return nil
}
