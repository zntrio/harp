// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ruleset

import (
	"reflect"
	"testing"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func TestValidate(t *testing.T) {
	type args struct {
		spec *bundlev1.RuleSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "invalid apiVersion",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kind",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "nil meta",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "RuleSet",
				},
			},
			wantErr: true,
		},
		{
			name: "meta name not defined",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "RuleSet",
					Meta:       &bundlev1.RuleSetMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "nil spec",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "RuleSet",
					Meta:       &bundlev1.RuleSetMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "no action patch",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "RuleSet",
					Meta:       &bundlev1.RuleSetMeta{},
					Spec:       &bundlev1.RuleSetSpec{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.spec); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	type args struct {
		spec *bundlev1.RuleSet
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				spec: &bundlev1.RuleSet{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "RuleSet",
					Meta:       &bundlev1.RuleSetMeta{},
					Spec:       &bundlev1.RuleSetSpec{},
				},
			},
			wantErr: false,
			want:    "yM_TR6rMWW7BGA1Ms-U3WK6E4Xax5qRBMjK4VQLyZmQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Checksum(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
