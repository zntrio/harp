// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/jmespath/go-jmespath"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func Test_matchJMESPath_IsSatisfiedBy(t *testing.T) {
	type fields struct {
		exp *jmespath.JMESPath
	}
	type args struct {
		object interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "nil",
			want: false,
		},
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "not supported type",
			fields: fields{
				exp: jmespath.MustCompile("true"),
			},
			args: args{
				object: struct{}{},
			},
			want: false,
		},
		{
			name: "supported type: empty object",
			fields: fields{
				exp: jmespath.MustCompile("true"),
			},
			args: args{
				object: &bundlev1.Package{},
			},
			want: false,
		},
		{
			name: "supported type: nil exp",
			args: args{
				object: &bundlev1.Package{},
			},
			want: false,
		},
		{
			name: "supported type: not matching",
			fields: fields{
				exp: jmespath.MustCompile("name=='test'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: false,
		},
		{
			name: "supported type: annotations query with nil",
			fields: fields{
				exp: jmespath.MustCompile("annotations.patched=='foo'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: false,
		},
		{
			name: "supported type: annotations not matching",
			fields: fields{
				exp: jmespath.MustCompile("annotations.patched=='foo'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
					Annotations: map[string]string{
						"patched": "true",
					},
				},
			},
			want: false,
		},
		{
			name: "supported type: annotations not matching with same type",
			fields: fields{
				exp: jmespath.MustCompile("annotations.patched=='true'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
					Annotations: map[string]string{
						"patched": "false",
					},
				},
			},
			want: false,
		},
		{
			name: "supported type: annotations matching",
			fields: fields{
				exp: jmespath.MustCompile("annotations.patched=='true'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
					Annotations: map[string]string{
						"patched": "true",
					},
				},
			},
			want: true,
		},
		{
			name: "supported type: name matching",
			fields: fields{
				exp: jmespath.MustCompile("name=='foo'"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &jmesPathMatcher{
				exp: tt.fields.exp,
			}
			if got := s.IsSatisfiedBy(tt.args.object); got != tt.want {
				t.Errorf("jmesPathMatcher.IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchJMESPath_IsSatisfiedBy_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var (
			expr *jmespath.JMESPath
		)

		f.Fuzz(&expr)

		// Execute
		s := &jmesPathMatcher{
			exp: expr,
		}
		s.IsSatisfiedBy(&bundlev1.Package{
			Name: "foo",
		})
	}
}
