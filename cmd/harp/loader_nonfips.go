// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build !fips

package main

import (
	// Register hash functions
	//nolint:gosec // For legacy compatibility
	_ "crypto/md5"
	//nolint:gosec // For legacy compatibility
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"

	// Register encryption transformers.
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/aead"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/age"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/branca"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/dae"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/fernet"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/jwe"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/paseto"
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/secretbox"
	_ "zntr.io/harp/v2/pkg/sdk/value/signature/jws"
	_ "zntr.io/harp/v2/pkg/sdk/value/signature/paseto"
	_ "zntr.io/harp/v2/pkg/sdk/value/signature/raw"
	_ "zntr.io/harp/v2/pkg/vault"

	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/sha3"
)
