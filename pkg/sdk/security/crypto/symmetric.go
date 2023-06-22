// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package crypto

import (
	"encoding/base64"
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/fernet/fernet-go"
	"github.com/pkg/errors"
	"zntr.io/harp/v2/build/fips"
)

// -----------------------------------------------------------------------------

// Key generates symmetric encryption keys according to given keyType.
func Key(keyType string) (string, error) {
	switch keyType {
	case "aes:128":
		key := memguard.NewBufferRandom(16).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "aes:192":
		key := memguard.NewBufferRandom(24).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "aes:256":
		key := memguard.NewBufferRandom(32).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "aes:siv":
		if fips.Enabled() {
			return "", errors.New("aes:siv key generation is disabled in FIPS Mode")
		}
		key := memguard.NewBufferRandom(64).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "secretbox":
		if fips.Enabled() {
			return "", errors.New("secretbox key generation is disabled in FIPS Mode")
		}
		key := memguard.NewBufferRandom(32).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "chacha20":
		if fips.Enabled() {
			return "", errors.New("chacha20 key generation is disabled in FIPS Mode")
		}
		key := memguard.NewBufferRandom(32).Bytes()
		return base64.StdEncoding.EncodeToString(key), nil
	case "fernet":
		// Generate a fernet key
		k := &fernet.Key{}
		if err := k.Generate(); err != nil {
			return "", err
		}
		return k.Encode(), nil
	default:
		return "", fmt.Errorf("invalid keytype (%s) [aes:128, aes:192, aes:256, aes:siv, secretbox, chacha20, fernet]", keyType)
	}
}
