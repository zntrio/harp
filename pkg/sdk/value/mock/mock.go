// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package mock

import (
	"context"

	"zntr.io/harp/v2/pkg/sdk/value"
)

func Transformer(err error) value.Transformer {
	return &mockedTransformer{
		err: err,
	}
}

type mockedTransformer struct {
	err error
}

func (m *mockedTransformer) To(ctx context.Context, input []byte) ([]byte, error) {
	return input, m.err
}

func (m *mockedTransformer) From(ctx context.Context, input []byte) ([]byte, error) {
	return input, m.err
}
