// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package convert

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

func Test_PBtoYAML(t *testing.T) {
	spec := &bundlev1.Patch{
		ApiVersion: "harp.zntr.io/v2",
		Kind:       "BundlePatch",
		Meta: &bundlev1.PatchMeta{
			Name: "test-patch",
		},
		Spec: &bundlev1.PatchSpec{
			Rules: []*bundlev1.PatchRule{
				{
					Package:  &bundlev1.PatchPackage{},
					Selector: &bundlev1.PatchSelector{},
				},
			},
		},
	}

	expectedOutput := []byte("apiVersion: harp.zntr.io/v2\nkind: BundlePatch\nmeta:\n  name: test-patch\nspec:\n  rules:\n  - package: {}\n    selector: {}\n")

	out, err := PBtoYAML(spec)
	if err != nil {
		t.Error(err)
	}

	if report := cmp.Diff(string(out), string(expectedOutput)); report != "" {
		t.Errorf("unexpected conversion output:\n%v", report)
	}
}

func Test_convertMapStringInterface(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: false,
		},
		{
			name: "map[interface{}]interface{}",
			args: args{
				val: map[interface{}]interface{}{
					"abc":  1234,
					"true": 12.56,
				},
			},
			wantErr: false,
			want: map[string]interface{}{
				"abc":  1234,
				"true": 12.56,
			},
		},
		{
			name: "[]interface{}",
			args: args{
				val: []interface{}{
					"abc", "true",
				},
			},
			wantErr: false,
			want:    []interface{}{"abc", "true"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertMapStringInterface(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertMapStringInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertMapStringInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}
