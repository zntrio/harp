// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/awnumar/memguard"
	"google.golang.org/protobuf/proto"

	containerv1 "zntr.io/harp/v2/api/gen/go/harp/container/v1"
	"zntr.io/harp/v2/pkg/sdk/types"

	"golang.org/x/crypto/nacl/secretbox"
)

// Unseal a sealed container with the given identity.
func (a *adapter) Unseal(container *containerv1.Container, identity *memguard.LockedBuffer) (*containerv1.Container, error) {
	return a.unseal(container, identity, nil)
}

// Unseal a sealed container with the given identity and the given preshared key.
func (a *adapter) UnsealWithPSK(container *containerv1.Container, identity, preSharedKey *memguard.LockedBuffer) (*containerv1.Container, error) {
	return a.unseal(container, identity, preSharedKey)
}

//nolint:gocyclo,funlen // To refactor
func (a *adapter) unseal(container *containerv1.Container, identity, preSharedKey *memguard.LockedBuffer) (*containerv1.Container, error) {
	// Check parameters
	if types.IsNil(container) {
		return nil, fmt.Errorf("unable to process nil container")
	}
	if types.IsNil(container.Headers) {
		return nil, fmt.Errorf("unable to process nil container headers")
	}
	if identity == nil {
		return nil, fmt.Errorf("unable to process without container key")
	}

	// Check headers
	if container.Headers.ContentType != containerSealedContentType {
		return nil, fmt.Errorf("unable to unseal container")
	}

	// Check ephemeral container public encryption key
	if len(container.Headers.EncryptionPublicKey) != publicKeySize {
		return nil, fmt.Errorf("invalid container public size")
	}
	var publicKey [publicKeySize]byte
	copy(publicKey[:], container.Headers.EncryptionPublicKey[:publicKeySize])

	// Decode private key
	privRaw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(identity.String(), PrivateKeyPrefix))
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	if len(privRaw) != privateKeySize {
		return nil, fmt.Errorf("invalid identity private key length")
	}
	var pk [privateKeySize]byte
	copy(pk[:], privRaw[:privateKeySize])

	// Compute preshared key
	var psk *[preSharedKeySize]byte
	if preSharedKey != nil {
		psk, err = pskStretch(preSharedKey.Bytes(), publicKey[:])
		if err != nil {
			return nil, fmt.Errorf("unable to stretch preshared key: %w", err)
		}
	}

	// Precompute identifier
	derivedKey, err := deriveSharedKeyFromRecipient(&publicKey, &pk, psk)
	if err != nil {
		return nil, fmt.Errorf("unable to compute recipient shared key: %w", err)
	}

	// Try recipients
	payloadKey, err := tryRecipientKeys(derivedKey, container.Headers.Recipients, psk)
	if err != nil {
		return nil, fmt.Errorf("error occurred during recipient key tests: %w", err)
	}

	// Check private key
	if len(payloadKey) != encryptionKeySize {
		return nil, fmt.Errorf("invalid encryption key size")
	}
	var encryptionKey [encryptionKeySize]byte
	copy(encryptionKey[:], payloadKey[:encryptionKeySize])

	// Prepare sig nonce
	var pubSigNonce [nonceSize]byte
	copy(pubSigNonce[:], "harp_container_psigk_box")

	// Decrypt signing public key
	containerSignKeyRaw, ok := secretbox.Open(nil, container.Headers.ContainerBox, &pubSigNonce, &encryptionKey)
	if !ok {
		return nil, fmt.Errorf("invalid container key")
	}
	if len(containerSignKeyRaw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid signature key size")
	}

	// Compute headers hash
	headerHash, err := computeHeaderHash(container.Headers)
	if err != nil {
		return nil, fmt.Errorf("unable to compute header hash: %w", err)
	}

	// Extract payload nonce
	var payloadNonce [nonceSize]byte
	copy(payloadNonce[:], headerHash[:nonceSize])

	// Decrypt payload
	payloadRaw, ok := secretbox.Open(nil, container.Raw, &payloadNonce, &encryptionKey)
	if !ok || len(payloadRaw) < signatureSize {
		return nil, fmt.Errorf("invalid ciphered content")
	}

	// Extract signature / content
	detachedSig := payloadRaw[:signatureSize]
	content := payloadRaw[signatureSize:]

	// Prepare protected content
	protectedHash := computeProtectedHash(headerHash, content)

	// Validate signature
	if !ed25519.Verify(containerSignKeyRaw, protectedHash, detachedSig) {
		return nil, fmt.Errorf("invalid container signature")
	}

	// Unmarshal inner container
	out := &containerv1.Container{}
	if err := proto.Unmarshal(content, out); err != nil {
		return nil, fmt.Errorf("unable to unpack inner content: %w", err)
	}

	// No error
	return out, nil
}
