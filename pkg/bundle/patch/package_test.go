// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package patch

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

var (
	opt = cmp.FilterPath(
		func(p cmp.Path) bool {
			// Remove ignoring of the fields below once go-cmp is able to ignore generated fields.
			// See https://github.com/google/go-cmp/issues/153
			ignoreXXXCache :=
				p.String() == "XXX_sizecache" ||
					p.String() == "Packages.XXX_sizecache" ||
					p.String() == "Packages.Secrets.XXX_sizecache" ||
					p.String() == "Packages.Secrets.Data.XXX_sizecache"
			return ignoreXXXCache
		}, cmp.Ignore())

	ignoreOpts = []cmp.Option{
		cmpopts.IgnoreUnexported(bundlev1.Bundle{}),
		cmpopts.IgnoreUnexported(bundlev1.Package{}),
		cmpopts.IgnoreUnexported(bundlev1.SecretChain{}),
		cmpopts.IgnoreUnexported(bundlev1.KV{}),
		opt,
	}
)

func TestValidate(t *testing.T) {
	type args struct {
		spec *bundlev1.Patch
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
				spec: &bundlev1.Patch{
					ApiVersion: "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid kind",
			args: args{
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "nil meta",
			args: args{
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundlePatch",
				},
			},
			wantErr: true,
		},
		{
			name: "meta name not defined",
			args: args{
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundlePatch",
					Meta:       &bundlev1.PatchMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "nil spec",
			args: args{
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundlePatch",
					Meta:       &bundlev1.PatchMeta{},
				},
			},
			wantErr: true,
		},
		{
			name: "no action patch",
			args: args{
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundlePatch",
					Meta:       &bundlev1.PatchMeta{},
					Spec:       &bundlev1.PatchSpec{},
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
		spec *bundlev1.Patch
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
				spec: &bundlev1.Patch{
					ApiVersion: "harp.zntr.io/v2",
					Kind:       "BundlePatch",
					Meta:       &bundlev1.PatchMeta{},
					Spec:       &bundlev1.PatchSpec{},
				},
			},
			wantErr: false,
			want:    "-1TA5k2pJzUSVu2HechzEWeS1Wx-H775AqfdOwUV2Ow",
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

func TestApply_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 500; i++ {
		f := fuzz.New()

		// Prepare arguments
		values := map[string]interface{}{}
		spec := &bundlev1.Patch{
			ApiVersion: "harp.zntr.io/v2",
			Kind:       "BundlePatch",
			Meta: &bundlev1.PatchMeta{
				Name: "test-patch",
			},
			Spec: &bundlev1.PatchSpec{
				Executor: &bundlev1.PatchExecutor{},
				Rules: []*bundlev1.PatchRule{
					{
						Package:  &bundlev1.PatchPackage{},
						Selector: &bundlev1.PatchSelector{},
					},
				},
			},
		}
		file := bundlev1.Bundle{
			Packages: []*bundlev1.Package{
				{
					Name: "foo",
					Secrets: &bundlev1.SecretChain{
						Data: []*bundlev1.KV{
							{
								Key:   "k1",
								Value: []byte("v1"),
							},
						},
					},
				},
			},
		}

		f.Fuzz(&spec)
		f.Fuzz(&file)

		// Execute
		Apply(context.Background(), spec, &file, values)
	}
}

func mustLoadPatch(filePath string) *bundlev1.Patch {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	p, err := YAML(f)
	if err != nil {
		panic(err)
	}

	return p
}

func TestApply(t *testing.T) {
	type args struct {
		spec    *bundlev1.Patch
		b       *bundlev1.Bundle
		values  map[string]interface{}
		options []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *bundlev1.Bundle
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "empty bundle",
			args: args{
				spec:   mustLoadPatch("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				b:      &bundlev1.Bundle{},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{},
			},
		},
		{
			name: "modifiable bundle",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-1.yaml",
						},
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Annotations: map[string]string{
							"patched":             "true",
							"secret-path-cleaner": "true",
						},
						Name: "application/component-1",
					},
					{
						Annotations: map[string]string{
							"patched":             "true",
							"secret-path-cleaner": "true",
						},
						Name: "application/component-2",
					},
				},
			},
		},
		{
			name: "duplicate package paths",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-1.yaml",
						},
						{
							Name: "secrets/application/component-1.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Annotations: map[string]string{
							"patched":             "true",
							"secret-path-cleaner": "true",
						},
						Name: "application/component-1",
					},
					{
						Annotations: map[string]string{
							"patched":             "true",
							"secret-path-cleaner": "true",
						},
						Name: "application/component-1",
					},
				},
			},
		},
		{
			name: "remove package",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/remove-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "application/to-be-removed",
						},
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "remove package with rego",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/rego-remove-packages.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "application/to-be-removed",
							Labels: map[string]string{
								"to-remove": "true",
							},
						},
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "remove secrets with secret matcher",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/remove-secrets.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "application/to-be-removed",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key: "USER",
									},
								},
							},
						},
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/to-be-removed",
						Annotations: map[string]string{
							"secret-remover": "true",
							"patched":        "true",
						},
						Secrets: &bundlev1.SecretChain{
							Data: []*bundlev1.KV{},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "add package",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/add-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/another-created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "application/created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "add package - stop at 1",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/add-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
				options: []OptionFunc{
					WithStopAtRuleIndex(1),
				},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "add package - stop at 'another-package'",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/add-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
				options: []OptionFunc{
					WithStopAtRuleID("another-package"),
				},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "add package - ignore 0",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/add-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
				options: []OptionFunc{
					WithIgnoreRuleIndexes(0),
				},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/another-created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
		{
			name: "add package - ignore id",
			args: args{
				spec: mustLoadPatch("../../../test/fixtures/patch/valid/add-package.yaml"),
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secrets/application/component-2.yaml",
						},
					},
				},
				values: map[string]interface{}{},
				options: []OptionFunc{
					WithIgnoreRuleIDs("another-package"),
				},
			},
			wantErr: false,
			want: &bundlev1.Bundle{
				Packages: []*bundlev1.Package{
					{
						Name: "application/created-package",
						Annotations: map[string]string{
							"package-creator":                       "true",
							"patched":                               "true",
							"secret-service.elstc.co/encryptionKey": "UcbPlrEJ9jZEQX06n8oMln_mCl3EU2zl2ZVc-obb7Dw=",
						},
						Secrets: &bundlev1.SecretChain{
							Annotations: map[string]string{
								"secret-service.elstc.co/encryptionKey": "DrZ-0yEA18iS7A4xaR_pd-relh9KMtTw2q11nBEJykg=",
							},
							Data: []*bundlev1.KV{
								{
									Key:   "key",
									Type:  "string",
									Value: []byte("0\n\x02\x01\x01\x13\x05value"),
								},
							},
						},
					},
					{
						Name: "secrets/application/component-2.yaml",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Apply(context.Background(), tt.args.spec, tt.args.b, tt.args.values, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, ignoreOpts...); diff != "" {
				t.Errorf("%q. Patch.Apply():\n-got/+want\ndiff %s", tt.name, diff)
			}
		})
	}
}
