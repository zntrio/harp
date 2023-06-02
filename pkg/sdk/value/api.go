// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package value

import (
	"context"

	// For tests.
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination test/mock/transformer.gen.go -package mock zntr.io/harp/v2/pkg/sdk/value Transformer

// Transformer declares value transformer contract.
type Transformer interface {
	To(ctx context.Context, input []byte) ([]byte, error)
	From(ctx context.Context, input []byte) ([]byte, error)
}
