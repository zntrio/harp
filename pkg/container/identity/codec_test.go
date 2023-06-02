// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT
package identity

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/container/identity/key"
)

var (
	v1SecurityIdentity = []byte(`{"@apiVersion":"harp.zntr.io/v2","@kind":"ContainerIdentity","@timestamp":"2023-06-02T09:10:16.824302Z","@description":"security","public":"v1.ipk.AOMqsGUl9zS0tY8aeX_L2W52Qdj7MOOD-Vukcv7I_AA","private":{"content":"ZXlKaGJHY2lPaUpRUWtWVE1pMUlVelV4TWl0Qk1qVTJTMWNpTENKbGJtTWlPaUpCTWpVMlIwTk5JaXdpY0RKaklqbzFNREF3TURFc0luQXljeUk2SWxOWlZHc3RSM0JOYjBSRE9YWlNWMHBXWld3NVdXY2lmUS5SdWRmWkdvbGpReGVPTnBBdVUxYVdtMDdqUE1VckxCeGc0aHJoUVVTTkhRa2FjLUYtR3o4LUEuYkR3WWJJU2ZkallTY0pvdi5FQTE3b1VWT2VNRS10a3hFaGVtNS1tY0Jya1NiZEVZdERabUN5cktIeTN0ZnA3Rm1PcG9wa3Z6dTBrdlNHQjY1VmVFMUZ1RGNFZWhHVVA1a0hNaVJWRU9ycl9IRGhzOXl6Zzk1dG5QcjRNTy1PeklOLVNwSTlSYlJkRTFMMk9xR0swWHlXNFdydGZ2emQwNDlWUFVvSWswNkh0dWlzZ21RaUFsazVQS2ZKbzNqZjdsdzh6Qm5scVItS2ZJYTQ5RXBjTEUtRzdiUFUwNk9tUnhFYTg0c0lZQWl2RnRYZnEtWHlPLWpSUS5SLVJGamZrY0g4eDRXaUYxb1htTnlR"},"signature":"sYka-wQ_iyb5fboOw2lGPu0Q_XwDEeacIHMuTSbRZCzNoQB6C8XNjapzCjr5npxLPt2EhiCMsQGh69yoRjGRDw"}`)
	publicOnly         = []byte(`{"@apiVersion":"harp.zntr.io/v2","@kind":"ContainerIdentity","@timestamp":"2023-06-02T09:10:16.824302Z","@description":"security","public":"v1.ipk.AOMqsGUl9zS0tY8aeX_L2W52Qdj7MOOD-Vukcv7I_AA"}`)
	v2SecurityIdentity = []byte(`{"@apiVersion":"harp.zntr.io/v2","@kind":"ContainerIdentity","@timestamp":"2023-06-02T09:11:38.297756Z","@description":"security","public":"v2.ipk.AyoIQyBpS422KbgIzx0GkhQQ_PP-JRyIOeOadgQRgQrZ8fsGYrtUql5Tktyr3rEEOQ","private":{"content":"ZXlKaGJHY2lPaUpRUWtWVE1pMUlVelV4TWl0Qk1qVTJTMWNpTENKbGJtTWlPaUpCTWpVMlIwTk5JaXdpY0RKaklqbzFNREF3TURFc0luQXljeUk2SWtvMGJFVTJNRkptZUVGS2FtSjBNMmxVYUVKQmFsRWlmUS4zazlySWh2VTB3UURGb0lMbEtYTjYtS0pLQjRXR0Y5U2JnZUJ0WXlfcUg0M0pSbnlESzQ2Y2cuOHJvalV6QzU0blVpMjFQdi5PVzg3MnE1UEttRjZXRU9QaEhsQWVhTWFueDJpMXZhUHUya0J0eDh3U3dOLXA4R3owTGZURmpBQkZFSHg4ZThqMG5xZHh1UmZLZ05LNGZVY1NVZ0xKMTVGdG40ZV9VRWpwZE9qaktzdnlPWi1mVFZVV2RIMzZabTl6Y0JwLVVsWFVPNUNXU21sTXh1a1R3MEVwYW5EUXVVZEQza1B1aFl0UURVOFZxVWxzVFpKa1VCZ0xkTkRiNTIzWTRfYmMtRkZybWRRc2ZZZDBqeHczUjVvUVBjVzFfTXhpcThZWG1Ob3BhUjBaM1RNWXNLUjdoOWdraVZVcUJQa1FDUm9CNUgxR3ZWNGdBZDg4c2dXSjBNcVkwaEZhU29xWHEwR2Z5dWd0Q2dJWlF1VlloNlBlT0FoSjNmdVhZYVJ5SzROaU9ZLjdJN3RtVkNXTzVSTnpnZHJHYm5HdVE"},"signature":"OPQzNzJ1E4hHt5dIlEBjk8i3qfxUw-5iSt-IJlE1UDFsbKFdwRfsizpvEC2hvdI2Hf52fVCfLTLhzxGVckxg_zjBzBY0MCki0aU9zFLcoRWrfx_cASGyFFNC4PYZ53NM"}`)
)

func TestCodec_New(t *testing.T) {
	t.Run("invalid description", func(t *testing.T) {
		id, pub, err := New(rand.Reader, "é", key.Ed25519)
		assert.Error(t, err)
		assert.Nil(t, pub)
		assert.Nil(t, id)
	})

	t.Run("ed25519 - invalid random source", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer(nil), "test", key.Ed25519)
		assert.Error(t, err)
		assert.Nil(t, pub)
		assert.Nil(t, id)
	})

	t.Run("p384 - invalid random source", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer(nil), "test", key.P384)
		assert.Error(t, err)
		assert.Nil(t, pub)
		assert.Nil(t, id)
	})

	t.Run("legacy - invalid random source", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer(nil), "test", key.Legacy)
		assert.Error(t, err)
		assert.Nil(t, pub)
		assert.Nil(t, id)
	})

	t.Run("valid - ed25519", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer([]byte("deterministic-random-source-for-test-0001")), "security", key.Ed25519)
		assert.NoError(t, err)
		assert.NotNil(t, pub)
		assert.NotNil(t, id)
		assert.Equal(t, "harp.zntr.io/v2", id.APIVersion)
		assert.Equal(t, "security", id.Description)
		assert.Equal(t, "ContainerIdentity", id.Kind)
		assert.Equal(t, "v1.ipk.2BdsL_FTiaLRwyYwlA2urcZ8TLDdisbzBSEp-LUuHos", id.Public)
		assert.Nil(t, id.Private)
		assert.False(t, id.HasPrivateKey())
	})

	t.Run("valid - p-384", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer([]byte("deterministic-random-source-for-test-0001-1ioQiLEbVCm1Y7NfWCf6oNWoV6p5E4spJgRXKQHdV44XcNvqywMnIYYcL8qZ4Wk")), "security", key.P384)
		assert.NoError(t, err)
		assert.NotNil(t, pub)
		assert.NotNil(t, id)
		assert.Equal(t, "harp.zntr.io/v2", id.APIVersion)
		assert.Equal(t, "security", id.Description)
		assert.Equal(t, "ContainerIdentity", id.Kind)
		assert.Equal(t, "v2.ipk.A0X20rlE8Pqp-YoMG8SNOop918AyfoSF_R9Z7MF5vP5nUoc_ZSRWauQR6cL4DqgrRA", id.Public)
		assert.Nil(t, id.Private)
		assert.False(t, id.HasPrivateKey())
	})

	t.Run("valid - legacy", func(t *testing.T) {
		id, pub, err := New(bytes.NewBuffer([]byte("deterministic-random-source-for-test-0001")), "security", key.Legacy)
		assert.NoError(t, err)
		assert.NotNil(t, pub)
		assert.NotNil(t, id)
		assert.Equal(t, "harp.zntr.io/v2", id.APIVersion)
		assert.Equal(t, "security", id.Description)
		assert.Equal(t, "ContainerIdentity", id.Kind)
		assert.Equal(t, "ZxTKWxgrG341_FxatkkfAxedMtfz1zJzAm6FUmitxHM", id.Public)
		assert.Nil(t, id.Private)
		assert.False(t, id.HasPrivateKey())
	})
}

func TestCodec_FromReader(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		id, err := FromReader(nil)
		assert.Error(t, err)
		assert.Nil(t, id)
	})

	t.Run("empty", func(t *testing.T) {
		id, err := FromReader(bytes.NewReader([]byte("{}")))
		assert.Error(t, err)
		assert.Nil(t, id)
	})

	t.Run("invalid json", func(t *testing.T) {
		id, err := FromReader(bytes.NewReader([]byte("{")))
		assert.Error(t, err)
		assert.Nil(t, id)
	})

	t.Run("public key only", func(t *testing.T) {
		id, err := FromReader(bytes.NewReader(publicOnly))
		assert.Error(t, err)
		assert.Nil(t, id)
	})

	t.Run("valid - v1", func(t *testing.T) {
		id, err := FromReader(bytes.NewReader(v1SecurityIdentity))
		assert.NoError(t, err)
		assert.NotNil(t, id)
	})

	t.Run("valid - v2", func(t *testing.T) {
		id, err := FromReader(bytes.NewReader(v2SecurityIdentity))
		assert.NoError(t, err)
		assert.NotNil(t, id)
	})
}
