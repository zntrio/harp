// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package operation

import (
	"context"
)

// Operation describes operation contract.
type Operation interface {
	Run(ctx context.Context) error
}

const (
	legacyBundleMetadataPrefix = "harp.zntr.io/v1/bundle"
)
