// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package tasks

import (
	"context"
	"io"
)

// Task describes task contract.
type Task interface {
	Run(ctx context.Context) error
}

// ReaderProvider describes io.Reader provider.
type ReaderProvider func(ctx context.Context) (io.Reader, error)

// WriterProvider describes io.Writer provider.
type WriterProvider func(ctx context.Context) (io.Writer, error)

// ReadSeekerProvider describes io.ReadSeeker provider.
type ReadSeekerProvider func(ctx context.Context) (io.ReadSeeker, error)
