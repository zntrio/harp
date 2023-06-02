// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v2

import (
	"crypto/elliptic"

	"zntr.io/harp/v2/pkg/container/seal"
)

const (
	SealVersion = 2
)

const (
	containerSealedContentType = "application/vnd.harp.v1.SealedContainer"
	seedSize                   = 32
	publicKeySize              = 49
	privateKeySize             = 48
	encryptionKeySize          = 32
	preSharedKeySize           = 64
	nonceSize                  = 16
	macSize                    = 48
	signatureSize              = 96
	messageLimit               = 64 * 1024 * 1024
)

var (
	encryptionCurve = elliptic.P384()
	signatureCurve  = elliptic.P384()
)

// -----------------------------------------------------------------------------

func New() seal.Strategy {
	return &adapter{}
}

type adapter struct{}
