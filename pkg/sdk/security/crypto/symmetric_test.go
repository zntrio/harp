// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package crypto

import "testing"

func TestKey(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			args:    "",
			wantErr: true,
		},
		{
			name:    "invalid",
			args:    "foo",
			wantErr: true,
		},
		{
			name:    "aes-128",
			args:    "aes:128",
			wantErr: false,
		},
		{
			name:    "aes-256",
			args:    "aes:256",
			wantErr: false,
		},
		{
			name:    "secretbox",
			args:    "secretbox",
			wantErr: false,
		},
		{
			name:    "fernet",
			args:    "fernet",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Key(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
