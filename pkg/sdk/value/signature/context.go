// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package signature

import (
	"context"
)

var (
	contextKeyDetachedSignature      = contextKey("detachedSignature")
	contextKeyInputHash              = contextKey("inputHash")
	contextKeyDeterministicSignature = contextKey("deterministic")
)

// -----------------------------------------------------------------------------

type contextKey string

func (c contextKey) String() string {
	return "zntr.io/harp/v2/pkg/sdk/value/signature/" + string(c)
}

// -----------------------------------------------------------------------------

func WithDetachedSignature(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, contextKeyDetachedSignature, value)
}

func IsDetached(ctx context.Context) bool {
	value, ok := ctx.Value(contextKeyDetachedSignature).(bool)
	if !ok {
		return false
	}
	return value
}

func WithInputPreHashed(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, contextKeyInputHash, value)
}

func IsInputPreHashed(ctx context.Context) bool {
	value, ok := ctx.Value(contextKeyInputHash).(bool)
	if !ok {
		return false
	}
	return value
}

func WithDetermisticSignature(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, contextKeyDeterministicSignature, value)
}

func IsDeterministic(ctx context.Context) bool {
	value, ok := ctx.Value(contextKeyDeterministicSignature).(bool)
	if !ok {
		return false
	}
	return value
}
