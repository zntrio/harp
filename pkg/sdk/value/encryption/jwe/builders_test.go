// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jwe

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Transformer_JWE_InvalidKey(t *testing.T) {
	keys := []string{
		"",
		"foo",
		"123456",
		"a128kw:abSOB6OHnFKCHIm60OXsA==",
		"a192kw:b4JtOwQLOk1-RWxXUh5eG54nbdBihLT",
		"a256kw:TkxS6qSV6DBjn29JmU2ieMPnuCZNn3JelI1CDNqAQ8=",
		"a128kw:",
		"a192kw:",
		"a256kw:",
		"a512kw:TkxS6qSV6DBjn29JmU2ieMPnuCZNn3JelI1CDNqAQ8=",
	}
	for _, k := range keys {
		key := k
		t.Run(fmt.Sprintf("key `%s`", key), func(t *testing.T) {
			underTest, err := FromKey(key)
			if err == nil {
				t.Fatalf("Transformer should raise an error with key `%s`", key)
			}
			if underTest != nil {
				t.Fatalf("Transformer instance should be nil")
			}
		})
	}
}

func Test_FromKey(t *testing.T) {
	keys := []string{
		"a128kw:abSOB6OHnFK1CHIm60OXsA==",
		"a192kw:b4JtOwQLOks1-RWxXUh5eG54nbdBihLT",
		"a256kw:TkxS6qSV6eDBjn29JmU2ieMPnuCZNn3JelI1CDNqAQ8=",
		"pbes2-hs256-a128kw:stalemate-parkway-hardened-jeep-shrink-dimmer-platter-pretense",
		"pbes2-hs384-a192kw:stalemate-parkway-hardened-jeep-shrink-dimmer-platter-pretense",
		"pbes2-hs512-a256kw:stalemate-parkway-hardened-jeep-shrink-dimmer-platter-pretense",
	}
	for _, k := range keys {
		key := k
		t.Run(fmt.Sprintf("key `%s`", key), func(t *testing.T) {
			underTest, err := FromKey(key)
			assert.NoError(t, err)
			assert.NotNil(t, underTest)

			// Try to encrypt
			ctx := context.Background()
			encrypted, err := underTest.To(ctx, []byte("cleartext"))
			assert.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			// Try to decrypt
			out, err := underTest.From(ctx, encrypted)
			assert.NoError(t, err)
			assert.Equal(t, []byte("cleartext"), out)
		})
	}
}

func Test_TransformerKey(t *testing.T) {
	k := TransformerKey(AES128_KW, "abSOB6OHnFK1CHIm60OXsA==")
	assert.Equal(t, "a128kw:abSOB6OHnFK1CHIm60OXsA==", k)
}

func Test_Transformer_InvalidAlgorithm(t *testing.T) {
	tr, err := Transformer(KeyAlgorithm("a512kw"), "abSOB6OHnFK1CHIm60OXsA==")
	assert.Error(t, err)
	assert.Nil(t, tr)
}
