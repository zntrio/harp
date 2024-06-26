// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package jwe

import (
	"context"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/go-jose/go-jose/v3"
)

func mustDecodeBase64(in string) []byte {
	out, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	return out
}

func Test_jweTransformer_To(t *testing.T) {
	type fields struct {
		key               interface{}
		keyAlgorithm      jose.KeyAlgorithm
		contentEncryption jose.ContentEncryption
	}
	type args struct {
		in0   context.Context
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "a128kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-00001"),
				keyAlgorithm:      jose.A128KW,
				contentEncryption: jose.A128GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
		{
			name: "a192kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-00001"),
				keyAlgorithm:      jose.A192KW,
				contentEncryption: jose.A128GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
		{
			name: "a256kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-00001"),
				keyAlgorithm:      jose.A256KW,
				contentEncryption: jose.A256GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
		{
			name: "pbes2-hs256-a128kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-0001"),
				keyAlgorithm:      jose.PBES2_HS256_A128KW,
				contentEncryption: jose.A128GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
		{
			name: "pbes2-hs384-a192kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-0001"),
				keyAlgorithm:      jose.PBES2_HS384_A192KW,
				contentEncryption: jose.A192GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
		{
			name: "pbes2-hs512-a256kw",
			fields: fields{
				key:               []byte("deterministic-key-for-test-0001"),
				keyAlgorithm:      jose.PBES2_HS512_A256KW,
				contentEncryption: jose.A256GCM,
			},
			args: args{
				input: []byte("cleartext message"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &jweTransformer{
				key:               tt.fields.key,
				keyAlgorithm:      tt.fields.keyAlgorithm,
				contentEncryption: tt.fields.contentEncryption,
			}
			_, err := d.To(tt.args.in0, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("jweTransformer.To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_jweTransformer_From(t *testing.T) {
	type fields struct {
		key               interface{}
		keyAlgorithm      jose.KeyAlgorithm
		contentEncryption jose.ContentEncryption
	}
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty",
			fields: fields{
				key: (""),
			},
			args: args{
				input: []byte{},
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				key:               mustDecodeBase64("abSOB6OHnFK1CHIm60OXsA=="),
				keyAlgorithm:      jose.A128KW,
				contentEncryption: jose.A128GCM,
			},
			args: args{
				input: []byte("eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4R0NNIn0.22-PjbJqsJ6TFVLhPwYJG3a0HZq0cAcf.zWKWg_GfycXIrVa9.6XvjKMvr2CjG.pcMO_ou5QqTa6u6PzDWFIg"),
			},
			wantErr: false,
			want:    []byte("cleartext"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &jweTransformer{
				key:               tt.fields.key,
				keyAlgorithm:      tt.fields.keyAlgorithm,
				contentEncryption: tt.fields.contentEncryption,
			}
			got, err := d.From(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("jweTransformer.From() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jweTransformer.From() = %v, want %v", got, tt.want)
			}
		})
	}
}
