// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"

	"github.com/awnumar/memguard"
	"google.golang.org/protobuf/proto"
	containerv1 "zntr.io/harp/v2/api/gen/go/harp/container/v1"
	"zntr.io/harp/v2/pkg/sdk/security/crypto/extra25519"
	"zntr.io/harp/v2/pkg/sdk/types"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

// Seal a secret container with identities.
func (a *adapter) Seal(rand io.Reader, container *containerv1.Container, encodedPeersPublicKey ...string) (*containerv1.Container, error) {
	return a.seal(rand, container, nil, encodedPeersPublicKey...)
}

// Seal a secret container with identities and preshared key.
func (a *adapter) SealWithPSK(rand io.Reader, container *containerv1.Container, psk *memguard.LockedBuffer, encodedPeersPublicKey ...string) (*containerv1.Container, error) {
	return a.seal(rand, container, psk, encodedPeersPublicKey...)
}

//nolint:funlen,gocyclo // To refactor
func (a *adapter) seal(rand io.Reader, container *containerv1.Container, preSharedKey *memguard.LockedBuffer, encodedPeerPublicKeys ...string) (*containerv1.Container, error) {
	// Check parameters
	if types.IsNil(container) {
		return nil, fmt.Errorf("unable to process nil container")
	}
	if types.IsNil(container.Headers) {
		return nil, fmt.Errorf("unable to process nil container headers")
	}
	if len(encodedPeerPublicKeys) == 0 {
		return nil, fmt.Errorf("unable to process empty public keys")
	}

	// Convert public keys
	peerPublicKeys, err := a.publicKeys(encodedPeerPublicKeys...)
	if err != nil {
		return nil, fmt.Errorf("unable to convert peer public keys: %w", err)
	}

	// Serialize protobuf payload
	content, err := proto.Marshal(container)
	if err != nil {
		return nil, fmt.Errorf("unable to encode container content: %w", err)
	}

	// Check cleartext message size.
	if len(content) > messageLimit {
		return nil, errors.New("unable to seal the container, container is too large")
	}

	// Generate payload encryption key
	var payloadKey [encryptionKeySize]byte
	if _, err = io.ReadFull(rand, payloadKey[:]); err != nil {
		return nil, fmt.Errorf("unable to generate payload key for encryption")
	}

	// Generate ephemeral signing key
	sigPub, sigPriv, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, fmt.Errorf("unable to generate signing keypair")
	}

	// Encrypt public signature key
	var pubSigNonce [nonceSize]byte
	copy(pubSigNonce[:], staticSignatureNonce)
	encryptedPubSig := secretbox.Seal(nil, sigPub, &pubSigNonce, &payloadKey)
	memguard.WipeBytes(pubSigNonce[:])

	// Generate ephemeral encryption key
	encPub, encPriv, err := box.GenerateKey(rand)
	if err != nil {
		return nil, fmt.Errorf("unable to generate ephemeral encryption keypair")
	}

	// Prepare sealed container
	containerHeaders := &containerv1.Header{
		ContentType:         containerSealedContentType,
		EncryptionPublicKey: encPub[:],
		ContainerBox:        encryptedPubSig,
		Recipients:          []*containerv1.Recipient{},
		SealVersion:         SealVersion,
	}

	// Compute preshared key
	var psk *[preSharedKeySize]byte
	if preSharedKey != nil {
		psk, err = pskStretch(preSharedKey.Bytes(), containerHeaders.EncryptionPublicKey)
		if err != nil {
			return nil, fmt.Errorf("unable to stretch preshared key: %w", err)
		}
	}

	// Process recipients
	for _, peerPublicKey := range peerPublicKeys {
		if types.IsNil(peerPublicKey) {
			// Ignore nil key
			continue
		}
		if extra25519.IsEdLowOrder(peerPublicKey[:]) {
			return nil, fmt.Errorf("unable to process with low order public key")
		}

		// Pack recipient using its public key
		r, errPack := packRecipient(rand, &payloadKey, encPriv, peerPublicKey, psk)
		if errPack != nil {
			return nil, fmt.Errorf("unable to pack container recipient (%X): %w", *peerPublicKey, err)
		}

		// Append to container
		containerHeaders.Recipients = append(containerHeaders.Recipients, r)
	}

	// Sanity check
	if len(containerHeaders.Recipients) == 0 {
		return nil, errors.New("unable to seal a container without recipients")
	}

	// Compute header hash
	headerHash, err := computeHeaderHash(containerHeaders)
	if err != nil {
		return nil, fmt.Errorf("unable to compute header hash: %w", err)
	}

	// Prepare protected content
	protectedHash := computeProtectedHash(headerHash, content)

	// Sign th protected content
	containerSig := ed25519.Sign(sigPriv, protectedHash)

	// Prepare encryption nonce form sigHash
	var sigNonce [nonceSize]byte
	copy(sigNonce[:], headerHash[:nonceSize])

	// No error
	return &containerv1.Container{
		Headers: containerHeaders,
		Raw:     secretbox.Seal(nil, append(containerSig, content...), &sigNonce, &payloadKey),
	}, nil
}
