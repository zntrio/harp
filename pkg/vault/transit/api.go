// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package transit

import "context"

// Encryptor describes encryption operations contract.
type Encryptor interface {
	Encrypt(ctx context.Context, cleartext []byte) ([]byte, error)
}

// Decryptor describes decryption operations contract.
type Decryptor interface {
	Decrypt(ctx context.Context, encrypted []byte) ([]byte, error)
}

// Service represents the Vault Transit backend operation service contract.
type Service interface {
	Encryptor
	Decryptor
}
