// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package envelope

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption/secretbox"
)

type testEnvelopeService struct{}

func (s *testEnvelopeService) Encrypt(_ context.Context, data []byte) ([]byte, error) {
	return []byte(base64.URLEncoding.EncodeToString(data)), nil
}

func (s *testEnvelopeService) Decrypt(_ context.Context, data []byte) ([]byte, error) {
	return base64.URLEncoding.DecodeString(string(data))
}

// -----------------------------------------------------------------------------

func Test_Envelope_From(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "Payload",
			input:   []byte("foo"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arm mock
			ctx := context.Background()
			envelopeService := &testEnvelopeService{}

			underTest, err := Transformer(envelopeService, secretbox.Transformer)
			if err != nil {
				t.Errorf("error during transformer initialization, error = %v", err)
				return
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
			if diff := cmp.Diff(got, testCase.input); diff != "" {
				t.Errorf("%q. Envelope.From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

func Test_Envelope_To_From(t *testing.T) {
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
		},
		{
			name:    "Payload",
			input:   []byte("foo"),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arm mock
			ctx := context.Background()
			envelopeService := &testEnvelopeService{}

			underTest, err := Transformer(envelopeService, secretbox.Transformer)
			if err != nil {
				t.Errorf("error during transformer initialization, error = %v", err)
				return
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

			clearText, err := underTest.From(ctx, got)
			if err != nil {
				t.Errorf("error during the Fernet.From() call, error = %v", err)
			}
			if diff := cmp.Diff(clearText, testCase.input); diff != "" {
				t.Errorf("%q. Envelope.To/From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

type testErrorEnvelopeService struct{}

func (s *testErrorEnvelopeService) Encrypt(_ context.Context, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("foo")
}

func (s *testErrorEnvelopeService) Decrypt(_ context.Context, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("foo")
}

func Test_Envelope_Service_Error(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "Payload",
			input:   []byte("foo"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arm mock
			ctx := context.Background()
			envelopeService := &testErrorEnvelopeService{}

			underTest, err := Transformer(envelopeService, secretbox.Transformer)
			if err != nil {
				t.Errorf("error during transformer initialization, error = %v", err)
				return
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

			clearText, err := underTest.From(ctx, got)
			if err != nil {
				t.Errorf("error during the Fernet.From() call, error = %v", err)
			}
			if diff := cmp.Diff(clearText, testCase.input); diff != "" {
				t.Errorf("%q. Envelope.To/From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

func Test_Envelope_Transformer_Error(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "Payload",
			input:   []byte("foo"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arm mock
			ctx := context.Background()
			envelopeService := &testEnvelopeService{}

			underTest, err := Transformer(envelopeService, func(string) (value.Transformer, error) {
				return nil, fmt.Errorf("foo")
			})
			if err != nil {
				t.Errorf("error during transformer initialization, error = %v", err)
				return
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

			clearText, err := underTest.From(ctx, got)
			if err != nil {
				t.Errorf("error during the Fernet.From() call, error = %v", err)
			}
			if diff := cmp.Diff(clearText, testCase.input); diff != "" {
				t.Errorf("%q. Envelope.To/From():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}
