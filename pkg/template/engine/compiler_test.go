// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

import (
	"testing"
)

func TestValue(t *testing.T) {
	type args struct {
		input           string
		templateContext Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{},
			wantErr: true,
			want:    "",
		},
		{
			name: "invalid template",
			args: args{
				input: "{{ Values.foo }}",
				templateContext: &context{
					values: Values{"bar": "foo"},
				},
			},
			wantErr: true,
		},
		{
			name: "no values match",
			args: args{
				input: "{{ .Values.foo }}",
				templateContext: &context{
					values: Values{"bar": "foo"},
				},
			},
			wantErr: false,
			want:    "<no value>",
		},
		{
			name: "no values match (strict)",
			args: args{
				input: "{{ .Values.foo }}",
				templateContext: &context{
					values:     Values{"bar": "foo"},
					strictMode: true,
				},
			},
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				input: "{{ .Values.foo }}",
				templateContext: &context{
					values: Values{"foo": "bar"},
				},
			},
			wantErr: false,
			want:    "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderContext(tt.args.templateContext, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
