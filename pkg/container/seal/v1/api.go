// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"crypto/ed25519"

	"zntr.io/harp/v2/pkg/container/seal"
)

const (
	SealVersion = 1
)

const (
	containerSealedContentType = "application/vnd.harp.v1.SealedContainer"
	publicKeySize              = 32
	privateKeySize             = 32
	encryptionKeySize          = 32
	keyIdentifierSize          = 32
	nonceSize                  = 24
	preSharedKeySize           = 64
	signatureSize              = ed25519.SignatureSize
	messageLimit               = 64 * 1024 * 1024

	staticSignatureNonce      = "harp_container_psigk_box"
	signatureDomainSeparation = "harp encrypted signature"
)

// -----------------------------------------------------------------------------

func New() seal.Strategy {
	return &adapter{}
}

type adapter struct{}
