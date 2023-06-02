// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package template

import (
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/template/visitor/secretbuilder"
	"zntr.io/harp/v2/pkg/template/engine"
)

func TestValidate(t *testing.T) {
	type args struct {
		spec *bundlev1.Template
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
				spec: &bundlev1.Template{
					ApiVersion: "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kind",
			args: args{
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "nil meta",
			args: args{
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundleTemplate",
				},
			},
			wantErr: true,
		},
		{
			name: "meta name not defined",
			args: args{
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundleTemplate",
					Meta:       &bundlev1.TemplateMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "nil spec",
			args: args{
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundleTemplate",
					Meta:       &bundlev1.TemplateMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "no action template",
			args: args{
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundleTemplate",
					Meta:       &bundlev1.TemplateMeta{},
					Spec:       &bundlev1.TemplateSpec{},
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
		spec *bundlev1.Template
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
				spec: &bundlev1.Template{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundleTemplate",
					Meta:       &bundlev1.TemplateMeta{},
					Spec:       &bundlev1.TemplateSpec{},
				},
			},
			wantErr: false,
			want:    "3ipnWuWHabucGE3J1xGS0W5GxWtwTOeiI8kb5PucRkY",
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

func TestExecute_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		spec := &bundlev1.Template{
			ApiVersion: "harp.zntr.io/v2",
			Kind:       "BundleTemplate",
			Meta:       &bundlev1.TemplateMeta{},
			Spec:       &bundlev1.TemplateSpec{},
		}

		// Fuzz input
		f.Fuzz(&spec.Spec.Namespaces)

		// Initialize a bundle creator
		var b *bundlev1.Bundle
		v := secretbuilder.New(b, engine.NewContext())

		// Execute
		Execute(spec, v)
	}
}
