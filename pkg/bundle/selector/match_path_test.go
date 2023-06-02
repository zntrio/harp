// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"regexp"
	"testing"

	"github.com/gobwas/glob"
	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func Test_matchPath_IsSatisfiedBy(t *testing.T) {
	type fields struct {
		strict string
		regex  *regexp.Regexp
		g      glob.Glob
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
			args: args{
				object: struct{}{},
			},
			want: false,
		},
		{
			name: "supported type: path nil",
			args: args{
				object: &bundlev1.Package{},
			},
			want: false,
		},
		{
			name: "supported type: path empty",
			args: args{
				object: &bundlev1.Package{
					Name: "",
				},
			},
			want: false,
		},
		{
			name: "supported type: strict mode not match",
			fields: fields{
				strict: "bar",
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: false,
		},
		{
			name: "supported type: strict mode with match",
			fields: fields{
				strict: "foo",
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: true,
		},
		{
			name: "supported type: regexp mode not match",
			fields: fields{
				regex: regexp.MustCompile("bar"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: false,
		},
		{
			name: "supported type: regexp mode with match",
			fields: fields{
				regex: regexp.MustCompile("foo"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: true,
		},
		{
			name: "supported type: glob mode not match",
			fields: fields{
				g: glob.MustCompile("test"),
			},
			args: args{
				object: &bundlev1.Package{
					Name: "foo",
				},
			},
			want: false,
		},
		{
			name: "supported type: glob mode with match",
			fields: fields{
				g: glob.MustCompile("foo"),
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
			s := &matchPath{
				strict: tt.fields.strict,
				regex:  tt.fields.regex,
				g:      tt.fields.g,
			}
			if got := s.IsSatisfiedBy(tt.args.object); got != tt.want {
				t.Errorf("matchPath.IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchPath_IsSatisfiedBy_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var (
			strict string
			re     *regexp.Regexp
		)

		f.Fuzz(&strict)
		// f.Fuzz(&re)

		// Execute
		s := &matchPath{
			strict: strict,
			regex:  re,
		}
		s.IsSatisfiedBy(&bundlev1.Package{
			Name: "foo",
		})
	}
}
