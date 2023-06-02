// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT
package v2

import (
	"bytes"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/container/seal"
)

func TestGenerateKey(t *testing.T) {
	adapter := New()

	t.Run("deterministic", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(
			seal.WithDeterministicKey(memguard.NewBufferFromBytes([]byte("deterministic-seed-for-test-00001")), "Release 64"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "v2.ck.QwUEpYFxXpwFGrHQbHXGH0k4w_g9iDw38d67f9YHZwhvmEyE0R3McDMYr260lNck", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v2.sk.AuSjVpMZben6n9fXiaDj8bMjSvhcZ9n7c82VOt7v9_UBzZJaMLamkQUFAVp_9frpAg", pub)
	})

	t.Run("deterministic - same key with different target", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(
			seal.WithDeterministicKey(memguard.NewBufferFromBytes([]byte("deterministic-seed-for-test-00001")), "Release 65"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "v2.ck.2pWmwDtEjYAsLMR-7es_p3IvyYNrc3qSo5KbqrYmbCq5COcquwpr3SDnOmJrrbDp", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v2.sk.AwzwXF1XaZVry-pppsQ1ovSIMLtix-Nhq8NkBDEp46ulrHuY2onMg2_VusdD5D2YXg", pub)
	})

	t.Run("master key too short", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(
			seal.WithDeterministicKey(memguard.NewBufferFromBytes([]byte("too-short-masterkey")), "Release 64"),
		)
		assert.Error(t, err)
		assert.Empty(t, pk)
		assert.Empty(t, pub)
	})

	t.Run("default with given random source", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(seal.WithRandom(bytes.NewReader([]byte("UlLYMVJzTrAv0KYbl2KqCo9fnsyPLu9YNAO5iUsABeYMmkKe2TnSp8JLD9zThZk"))))
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "v2.ck.VHJBdjBLWWJsMktxQ285ZoFXc5G4HY_0qSMZAibGlchUmqt915byglIOGeel-5X5", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v2.sk.A0V1xCxGNtVAE9EVhaKi-pIADhd1in8xV_FI5Y0oHSHLAkew9gDAqiALSd6VgvBCbQ", pub)
	})

	t.Run("default", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey()
		assert.NoError(t, err)
		assert.NotEmpty(t, pk)
		assert.NotEmpty(t, pub)
	})
}
