// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package operation

import (
	"testing"
)

func Test_extractVersion(t *testing.T) {
	type args struct {
		packagePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   uint32
		wantErr bool
	}{
		{
			name:    "blank",
			wantErr: true,
		},
		{
			name: "no version",
			args: args{
				packagePath: "app/test",
			},
			wantErr: false,
			want:    "app/test",
			want1:   0,
		},
		{
			name: "with version",
			args: args{
				packagePath: "app/test?version=14",
			},
			wantErr: false,
			want:    "app/test",
			want1:   14,
		},
		{
			name: "with invalid version",
			args: args{
				packagePath: "app/test?version=azerty",
			},
			wantErr: true,
		},
		{
			name: "with invalid path",
			args: args{
				packagePath: "\n\t",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := extractVersion(tt.args.packagePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("exporter.extractVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exporter.extractVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("exporter.extractVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
