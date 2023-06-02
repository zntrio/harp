// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package identity

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Identity_From(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: false,
			want:    []byte{},
		},
		{
			name:    "Something",
			input:   []byte("foo"),
			wantErr: false,
			want:    []byte("foo"),
		},
	}
	for _, tC := range testCases {
		testCase := tC
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			underTest := Transformer()

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
			if diff := cmp.Diff(got, testCase.input); diff != "" {
				t.Errorf("%q. Identity.From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

func Test_Identity_To(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: false,
			want:    []byte{},
		},
		{
			name:    "Something",
			input:   []byte("foo"),
			wantErr: false,
			want:    []byte("foo"),
		},
	}
	for _, tC := range testCases {
		testCase := tC
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			underTest := Transformer()

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
			if diff := cmp.Diff(got, testCase.input); diff != "" {
				t.Errorf("%q. Identity.To():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}
