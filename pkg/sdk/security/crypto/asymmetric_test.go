// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package crypto

import (
	"testing"
)

func TestKeypair(t *testing.T) {
	type testCase struct {
		name    string
		args    string
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "nil",
			args:    "",
			wantErr: true,
		},
		{
			name:    "invalid",
			args:    "azer",
			wantErr: true,
		},
	}
	expectedKeyTypes := []string{"rsa", "rsa:normal", "rsa:2048", "rsa:strong", "rsa:4096", "ec", "ec:normal", "ec:p256", "ec:high", "ec:p384", "ec:strong", "ec:p521", "ssh", "ed25519", "naclbox"}
	for _, kt := range expectedKeyTypes {
		tests = append(tests, testCase{
			name:    kt,
			args:    kt,
			wantErr: false,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Keypair(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keypair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("Keypair() = %v", got)
			}
		})
	}
}
