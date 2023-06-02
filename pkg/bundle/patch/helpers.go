// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package patch

import bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"

// WithAnnotations returns the given patch spec patch annotations state.
func WithAnnotations(p *bundlev1.Patch) bool {
	switch {
	case p == nil, p.Spec == nil, p.Spec.Executor == nil:
		return true
	default:
		return !p.Spec.Executor.DisableAnnotations
	}
}
