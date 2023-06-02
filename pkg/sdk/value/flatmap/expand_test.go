// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package flatmap

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"zntr.io/harp/v2/pkg/bundle"
)

func TestExpand(t *testing.T) {
	type args struct {
		m   bundle.KV
		key string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "empty",
		},
		{
			name: "map",
			args: args{
				m: bundle.KV{
					"app/database": bundle.KV{
						"user":     "test",
						"password": "password",
					},
				},
			},
			want: bundle.KV{
				"app": bundle.KV{
					"database": bundle.KV{
						"user":     "test",
						"password": "password",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expand(tt.args.m, tt.args.key)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Expand() = \n%s", diff)
			}
		})
	}
}
