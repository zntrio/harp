// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package paseto

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/sdk/value"
)

var p384PrivateJWK = []byte(`{
    "kty": "EC",
    "d": "7YcsmkNxmZdzGyb46ZeDb2I1yr-ja1iw9gspGjq7UDqQ6a61h_ES8c4uU__adkFV",
    "crv": "P-384",
    "x": "dWLSo6PTkL1G68bzTwY3zzrL_QX-pwvP9HUPpQGeSFmj20EWOtfvXXKDrCR0jnJD",
    "y": "lFvTFechH_KmbOEvycryCHy23Cm1qekJYAtn7T_TELpm7zsY290NYlvqDKesGeXx"
}`)

var ed25519PrivateJWK = []byte(`{
    "kty": "OKP",
    "d": "ytOw6kKTTVJUKCnX5HgmhsGguNFQ18ECIS2C-ujJv-s",
    "crv": "Ed25519",
    "x": "K5i0d37-eRk8-EPwo2bpcmM-HGmzLiqRtWnk7oR3FCs"
}`)

func TestTransformer(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    value.Transformer
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "invalid base64",
			args: args{
				key: "paseto:123456789%",
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid - v3",
			args: args{
				key: fmt.Sprintf("paseto:%s", base64.RawURLEncoding.EncodeToString(p384PrivateJWK)),
			},
			wantErr: false,
		},
		{
			name: "valid - v4",
			args: args{
				key: fmt.Sprintf("paseto:%s", base64.RawURLEncoding.EncodeToString(ed25519PrivateJWK)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transformer(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
