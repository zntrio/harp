// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"io"

	"github.com/awnumar/memguard"

	containerv1 "zntr.io/harp/v2/api/gen/go/harp/container/v1"
	"zntr.io/harp/v2/pkg/sdk/types"
)

// Seal a secret container with identities.
func (a *adapter) Seal(rand io.Reader, container *containerv1.Container, encodedPeersPublicKey ...string) (*containerv1.Container, error) {
	return a.seal(rand, container, nil, encodedPeersPublicKey...)
}

// Seal a secret container with identities and preshared key.
func (a *adapter) SealWithPSK(rand io.Reader, container *containerv1.Container, preSharedKey *memguard.LockedBuffer, encodedPeersPublicKey ...string) (*containerv1.Container, error) {
	return a.seal(rand, container, preSharedKey, encodedPeersPublicKey...)
}

func (a *adapter) seal(rand io.Reader, container *containerv1.Container, preSharedKey *memguard.LockedBuffer, encodedPeersPublicKey ...string) (*containerv1.Container, error) {
	// Check parameters
	if types.IsNil(container) {
		return nil, fmt.Errorf("unable to process nil container")
	}
	if types.IsNil(container.Headers) {
		return nil, fmt.Errorf("unable to process nil container headers")
	}
	if len(encodedPeersPublicKey) == 0 {
		return nil, fmt.Errorf("unable to process empty public keys")
	}

	// Convert public keys
	peersPublicKey, err := a.publicKeys(encodedPeersPublicKey...)
	if err != nil {
		return nil, fmt.Errorf("unable to convert peer public keys: %w", err)
	}

	// Generate encryption key
	payloadKey, err := generatedEncryptionKey(rand)
	if err != nil {
		return nil, fmt.Errorf("unable to generate encryption key: %w", err)
	}

	// Prepare signature identity
	sigPriv, encryptedPubSig, err := prepareSignature(rand, payloadKey)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare signature materials: %w", err)
	}

	// Generate ephemeral encryption key
	encPriv, err := ecdsa.GenerateKey(encryptionCurve, rand)
	if err != nil {
		return nil, fmt.Errorf("unable to generate ephemeral encryption keypair")
	}

	// Prepare sealed container
	containerHeaders := &containerv1.Header{
		ContentType:         containerSealedContentType,
		EncryptionPublicKey: elliptic.MarshalCompressed(encPriv.Curve, encPriv.PublicKey.X, encPriv.PublicKey.Y),
		ContainerBox:        encryptedPubSig,
		Recipients:          []*containerv1.Recipient{},
		SealVersion:         SealVersion,
	}

	// Compute preshared key
	var psk *[preSharedKeySize]byte
	if preSharedKey != nil {
		psk = pskStretch(preSharedKey.Bytes(), containerHeaders.EncryptionPublicKey)
	}

	// Process recipients
	for _, peerPublicKey := range peersPublicKey {
		// Ignore nil key
		if peerPublicKey == nil {
			continue
		}

		// Pack recipient using its public key
		r, errPack := packRecipient(rand, payloadKey, encPriv, peerPublicKey, psk)
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

	// Sign given container
	content, containerSig, err := signContainer(sigPriv, containerHeaders, container)
	if err != nil {
		return nil, fmt.Errorf("unable to sign container data: %w", err)
	}

	// Encrypt payload
	encryptedPayload, err := encrypt(rand, append(containerSig, content...), payloadKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt container data: %w", err)
	}

	// No error
	return &containerv1.Container{
		Headers: containerHeaders,
		Raw:     encryptedPayload,
	}, nil
}
