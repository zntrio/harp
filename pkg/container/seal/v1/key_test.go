// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT
package v1

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
		assert.Equal(t, "v1.ck.8B_H8o7_ygAD27fFbqhgq97hLeJb5Nh4v3xy0C9JYPg", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v1.sk.qKXPnUP6-2Bb_4nYnmxOXyCdN4IV3AR5HooB33N3g2E", pub)
	})

	t.Run("deterministic - same key with different target", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(
			seal.WithDeterministicKey(memguard.NewBufferFromBytes([]byte("deterministic-seed-for-test-00001")), "Release 65"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "v1.ck.RIdVmnxg69ZKXkd7YknoIfvsnyfOTi792AhwlAIcaJ8", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v1.sk.SLP3GYe7UT-ADwuS2Ak-UEFCKR3ddvMawbwlgUSDG3k", pub)
	})

	t.Run("master key too short", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(
			seal.WithDeterministicKey(memguard.NewBufferFromBytes([]byte("determini")), "Release 64"),
		)
		assert.Error(t, err)
		assert.Empty(t, pk)
		assert.Empty(t, pub)
	})

	t.Run("default with given random source", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey(seal.WithRandom(bytes.NewReader([]byte("deterministic-seed-for-test-00001"))))
		assert.NoError(t, err)
		assert.NotNil(t, pk)
		assert.Equal(t, "v1.ck.ZGV0ZXJtaW5pc3RpYy1zZWVkLWZvci10ZXN0LTAwMDA", pk)
		assert.NotNil(t, pub)
		assert.Equal(t, "v1.sk.sYp90gC29yKfUUtr50pMR4Faf7c3d4-YX4xZsbwAs10", pub)
	})

	t.Run("default", func(t *testing.T) {
		pub, pk, err := adapter.GenerateKey()
		assert.NoError(t, err)
		assert.NotEmpty(t, pk)
		assert.NotEmpty(t, pub)
	})
}
