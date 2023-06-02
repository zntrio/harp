// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package identity

import (
	"context"

	"zntr.io/harp/v2/pkg/sdk/value"
)

type identityTransformer struct{}

// Transformer returns a non-operation transformer.
func Transformer() value.Transformer {
	return &identityTransformer{}
}

// -----------------------------------------------------------------------------

func (t *identityTransformer) From(_ context.Context, in []byte) ([]byte, error) {
	return in, nil
}

func (t *identityTransformer) To(_ context.Context, in []byte) ([]byte, error) {
	return in, nil
}
