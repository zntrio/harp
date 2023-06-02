// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package envelope

import "context"

//go:generate mockgen -destination test/mock/service.gen.go -package mock zntr.io/harp/v2/pkg/sdk/value/encryption/envelope Service

// Service declares envelope encryption service contract.
type Service interface {
	Decrypt(ctx context.Context, encrypted []byte) ([]byte, error)
	Encrypt(ctx context.Context, cleartext []byte) ([]byte, error)
}
