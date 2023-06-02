// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbox

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	keyLength   = 32
	nonceLength = 24
)

func generateNonce() ([nonceLength]byte, error) {
	var nonce [nonceLength]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	return nonce, err
}

func encrypt(plaintext []byte, key [keyLength]byte) ([]byte, error) {
	nonce, err := generateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce")
	}
	return secretbox.Seal(nonce[:], plaintext, &nonce, &key), nil
}

func decrypt(ciphertext []byte, key [keyLength]byte) ([]byte, error) {
	var nonce [nonceLength]byte
	copy(nonce[:], ciphertext[:nonceLength])
	decrypted, ok := secretbox.Open(nil, ciphertext[nonceLength:], &nonce, &key)
	if !ok {
		return nil, errors.New("failed to decrypt given message")
	}
	return decrypted, nil
}
