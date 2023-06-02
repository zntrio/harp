// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !fips

package main

import (
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
