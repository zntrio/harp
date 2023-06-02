// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"context"
	"testing"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func Test_matchRego_IsSatisfiedBy(t *testing.T) {
	type fields struct {
		policy string
	}
	type args struct {
		object interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "not supported type",
			fields: fields{
				policy: "package harp\ndefault matched = false\nmatched { input.labels }",
			},
			args: args{
				object: struct{}{},
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "supported type: empty object",
			fields: fields{
				policy: "package harp\ndefault matched = false\nmatched { input.labels }",
			},
			args: args{
				object: &bundlev1.Package{},
			},
			want: false,
		},
		{
			name: "supported type: invalid return",
			fields: fields{
				policy: "package harp\ndefault matched = \"\"",
			},
			args: args{
				object: &bundlev1.Package{},
			},
			want: false,
		},
		{
			name: "supported type: match",
			fields: fields{
				policy: "package harp\ndefault matched = false\nmatched { input.labels }",
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
					Labels: map[string]string{
						"patched": "true",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := MatchRego(context.Background(), tt.fields.policy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error got %v, expected %v", err, tt.wantErr)
				return
			}
			if s == nil {
				return
			}
			if got := s.IsSatisfiedBy(tt.args.object); got != tt.want {
				t.Errorf("regoMatcher.IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
