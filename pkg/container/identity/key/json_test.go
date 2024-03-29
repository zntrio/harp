// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	legacyPrivateKey = &JSONWebKey{
		Kty: "OKP",
		Crv: "X25519",
		X:   "ZxTKWxgrG341_FxatkkfAxedMtfz1zJzAm6FUmitxHM",
		D:   "ZGV0ZXJtaW5pc3RpYy1yYW5kb20tc291cmNlLWZvci0",
	}
	v1PrivateKey = &JSONWebKey{
		Kty: "OKP",
		Crv: "Ed25519",
		X:   "2BdsL_FTiaLRwyYwlA2urcZ8TLDdisbzBSEp-LUuHos",
		D:   "ZGV0ZXJtaW5pc3RpYy1yYW5kb20tc291cmNlLWZvci3YF2wv8VOJotHDJjCUDa6txnxMsN2KxvMFISn4tS4eiw",
	}
	v2PrivateKey = &JSONWebKey{
		Kty: "EC",
		Crv: "P-384",
		X:   "RfbSuUTw-qn5igwbxI06in3XwDJ-hIX9H1nswXm8_mdShz9lJFZq5BHpwvgOqCtE",
		Y:   "ag16lWruEPkhWChmZnO52ne1iyLGAEVNbyx38NPMOqNZzV7yP9ugrzCa7pCz8eBr",
		D:   "aXN0aWMtcmFuZG9tLXNvdYiXCnZ-xg0Te8QN3AId4n-bdBdDfhXJjz1OngEo78g8",
	}
)

func TestJSONWebKey_RecoveryKey(t *testing.T) {
	t.Run("D has invalid encoding", func(t *testing.T) {
		id, err := (&JSONWebKey{
			D: "é",
		}).RecoveryKey()
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("unhandled private key", func(t *testing.T) {
		id, err := (&JSONWebKey{
			Crv: "P-256",
		}).RecoveryKey()
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("valid - legacy", func(t *testing.T) {
		id, err := legacyPrivateKey.RecoveryKey()
		assert.NoError(t, err)
		assert.Equal(t, "ZGV0ZXJtaW5pc3RpYy1yYW5kb20tc291cmNlLWZvci0", id)
	})

	t.Run("valid - v1", func(t *testing.T) {
		id, err := v1PrivateKey.RecoveryKey()
		assert.NoError(t, err)
		assert.Equal(t, "v1.ck.6Of3g6qt-NPBzXSMNl4jPIZbrZIIwonT2pn7GCc4i3o", id)
	})

	t.Run("valid - v2", func(t *testing.T) {
		id, err := v2PrivateKey.RecoveryKey()
		assert.NoError(t, err)
		assert.Equal(t, "v2.ck.aXN0aWMtcmFuZG9tLXNvdYiXCnZ-xg0Te8QN3AId4n-bdBdDfhXJjz1OngEo78g8", id)
	})
}
