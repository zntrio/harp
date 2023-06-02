// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package fernet

import (
	"context"
	"fmt"
	"testing"

	"github.com/fernet/fernet-go"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func Test_Transformer_Fernet_InvalidKey(t *testing.T) {
	keys := []string{
		"",
		"foo",
		"123456",
	}
	for _, k := range keys {
		key := k
		t.Run(fmt.Sprintf("key `%s`", key), func(t *testing.T) {
			underTest, err := Transformer(key)
			if err == nil {
				t.Fatalf("Transformer should raise an error with key `%s`", key)
			}
			if underTest != nil {
				t.Fatalf("Transformer instance should be nil")
			}
		})
	}
}

func Test_Transformer_Fernet_From(t *testing.T) {
	// Generate a random fernet key
	k := &fernet.Key{}
	if err := k.Generate(); err != nil {
		t.Fatalf("unable to generate fernet key: %v", err)
	}

	// Prepare valid dataset
	plainText := []byte("cool-protected-data")
	encrypted, err := fernet.EncryptAndSign(plainText, k)
	if err != nil {
		t.Fatalf("unable to encrypt data with fernet key: %v", err)
	}

	// Prepare testcases
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Invalid encrypted payload",
			input:   []byte("bad-encryption-payload"),
			wantErr: true,
		},
		{
			name:    "Valid payload",
			input:   encrypted,
			wantErr: false,
			want:    plainText,
		},
	}

	// For each testcase
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Initialize mock
			ctx := context.Background()

			// Initialize transformer
			underTest, err := Transformer(k.Encode())
			if err != nil {
				t.Fatalf("unable to initialize transformer: %v", err)
			}

			// Do the call
			got, err := underTest.From(ctx, testCase.input)

			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if testCase.wantErr {
				return
			}
			if diff := cmp.Diff(got, testCase.want); diff != "" {
				t.Errorf("%q. Fernet.From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

func Test_Transformer_Fernet_To(t *testing.T) {
	// Generate a random fernet key
	k := &fernet.Key{}
	if err := k.Generate(); err != nil {
		t.Fatalf("unable to generate fernet key: %v", err)
	}

	t.Log(k.Encode())

	// Prepare valid dataset
	plainText := []byte("cool-protected-data")

	// Prepare testcases
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Valid payload",
			input:   plainText,
			wantErr: false,
		},
	}

	// For each testcase
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Initialize mock
			ctx := context.Background()

			// Initialize transformer
			underTest, err := Transformer(k.Encode())
			if err != nil {
				t.Fatalf("unable to initialize transformer: %v", err)
			}

			// Do the call
			got, err := underTest.To(ctx, testCase.input)

			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if testCase.wantErr {
				return
			}
			out, err := underTest.From(ctx, got)
			if err != nil {
				t.Errorf("error during the Fernet.From() call, error = %v", err)
			}
			if diff := cmp.Diff(out, testCase.input); diff != "" {
				t.Errorf("%q. Fernet.To():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}
