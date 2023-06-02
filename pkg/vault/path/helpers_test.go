// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package path

import "testing"

func TestSanitizePath(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "blank",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "whitespace prefixed",
			args: args{
				s: "  app/foo",
			},
			want: "app/foo",
		},
		{
			name: "whitespace suffixed",
			args: args{
				s: "app/foo   ",
			},
			want: "app/foo",
		},
		{
			name: "slash suffixed",
			args: args{
				s: "app/foo/",
			},
			want: "app/foo",
		},
		{
			name: "slash prefixed",
			args: args{
				s: "/app/foo",
			},
			want: "app/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizePath(tt.args.s); got != tt.want {
				t.Errorf("SanitizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
