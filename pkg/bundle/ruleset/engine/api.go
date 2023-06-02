// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

import (
	"context"
	"errors"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// ErrRuleNotValid is raised when a rule from a ruleset is false.
var ErrRuleNotValid = errors.New("rule is not valid")

// PackageLinter describes linter engine contract.
type PackageLinter interface {
	EvaluatePackage(ctx context.Context, p *bundlev1.Package) error
}
