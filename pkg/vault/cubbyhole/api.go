// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cubbyhole

import (
	"context"
	"io"
)

type Reader interface {
	Get(ctx context.Context, token string, w io.Writer) error
}

type Writer interface {
	Put(ctx context.Context, r io.Reader) (string, error)
}

type Service interface {
	Reader
	Writer
}
