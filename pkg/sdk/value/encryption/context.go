// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package encryption

import "context"

type contextKey string

func (c contextKey) String() string {
	return "zntr.io/harp/v2/pkg/sdk/value/encryption#" + string(c)
}

var (
	contextKeyNonce = contextKey("nonce")
	contextKeyAAD   = contextKey("aad")
)

func WithNonce(ctx context.Context, value []byte) context.Context {
	return context.WithValue(ctx, contextKeyNonce, value)
}

// Nonce gets the nonce value from the context.
func Nonce(ctx context.Context) ([]byte, bool) {
	nonce, ok := ctx.Value(contextKeyNonce).([]byte)
	return nonce, ok
}

func WithAdditionalData(ctx context.Context, value []byte) context.Context {
	return context.WithValue(ctx, contextKeyAAD, value)
}

// AdditionalData gets the aad value from the context.
func AdditionalData(ctx context.Context) ([]byte, bool) {
	aad, ok := ctx.Value(contextKeyAAD).([]byte)
	return aad, ok
}
