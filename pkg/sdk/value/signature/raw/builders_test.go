// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package raw

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/sdk/value"
)

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
				key: "raw:123456789%",
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid - v3",
			args: args{
				key: fmt.Sprintf("raw:%s", base64.RawURLEncoding.EncodeToString(p384PrivateJWK)),
			},
			wantErr: false,
		},
		{
			name: "valid - v4",
			args: args{
				key: fmt.Sprintf("raw:%s", base64.RawURLEncoding.EncodeToString(ed25519PrivateJWK)),
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
