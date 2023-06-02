// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbuilder

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	"zntr.io/harp/v2/pkg/template/engine"
)

func TestSuffix(t *testing.T) {
	type args struct {
		templateContext engine.Context
		ring            csov1.Ring
		secretPath      string
		item            *bundlev1.SecretSuffix
		data            interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{},
			wantErr: true,
		},
		{
			name: "invalid template function",
			args: args{
				templateContext: engine.NewContext(),
				ring:            csov1.RingInfra,
				secretPath:      "infra/aws/foo/us-east-1/rds/database/root_credentials",
				item: &bundlev1.SecretSuffix{
					Template: `{{ foo }}`,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			args: args{
				templateContext: engine.NewContext(),
				ring:            csov1.RingInfra,
				secretPath:      "infra/aws/foo/us-east-1/rds/database/root_credentials",
				item: &bundlev1.SecretSuffix{
					Template: `{"foo`,
				},
			},
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				templateContext: engine.NewContext(),
				ring:            csov1.RingInfra,
				secretPath:      "infra/aws/foo/us-east-1/rds/database/root_credentials",
				item: &bundlev1.SecretSuffix{
					Template: `{"foo":"123456"}`,
				},
			},
			wantErr: false,
			want: map[string]interface{}{
				"foo": "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderSuffix(tt.args.templateContext, tt.args.secretPath, tt.args.item, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Suffix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Suffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
