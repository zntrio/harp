// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package compare

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/secret"
)

func MustPack(value interface{}) []byte {
	out, err := secret.Pack(value)
	if err != nil {
		panic(err)
	}

	return out
}

func TestDiff(t *testing.T) {
	type args struct {
		src *bundlev1.Bundle
		dst *bundlev1.Bundle
	}
	tests := []struct {
		name    string
		args    args
		want    []DiffItem
		wantErr bool
	}{
		{
			name:    "src nil",
			wantErr: true,
		},
		{
			name: "dst nil",
			args: args{
				src: &bundlev1.Bundle{},
			},
			wantErr: true,
		},
		{
			name: "identic",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{},
				},
			},
			wantErr: false,
			want:    []DiffItem{},
		},
		{
			name: "new package",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key1",
										Value: MustPack("payload"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			want: []DiffItem{
				{Operation: Add, Type: "package", Path: "app/test"},
				{Operation: Add, Type: "secret", Path: "app/test#key1", Value: "payload"},
			},
		},
		{
			name: "package removed",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key1",
										Value: MustPack("payload"),
									},
								},
							},
						},
					},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{},
				},
			},
			wantErr: false,
			want: []DiffItem{
				{Operation: Remove, Type: "package", Path: "app/test"},
			},
		},
		{
			name: "secret added",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key1",
										Value: MustPack("payload"),
									},
								},
							},
						},
					},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key1",
										Value: MustPack("payload"),
									},
									{
										Key:   "key2",
										Value: MustPack("newpayload"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			want: []DiffItem{
				{Operation: Add, Type: "secret", Path: "app/test#key2", Value: "newpayload"},
			},
		},
		{
			name: "secret removed",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key1",
										Value: MustPack("payload"),
									},
									{
										Key:   "key2",
										Value: MustPack("newpayload"),
									},
								},
							},
						},
					},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key2",
										Value: MustPack("newpayload"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			want: []DiffItem{
				{Operation: Remove, Type: "secret", Path: "app/test#key1"},
			},
		},
		{
			name: "secret replaced",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key2",
										Value: MustPack("oldpayload"),
									},
								},
							},
						},
					},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key2",
										Value: MustPack("newpayload"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			want: []DiffItem{
				{Operation: Replace, Type: "secret", Path: "app/test#key2", Value: "newpayload"},
			},
		},
		{
			name: "no-op",
			args: args{
				src: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key2",
										Value: MustPack("oldpayload"),
									},
								},
							},
						},
					},
				},
				dst: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "app/test",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "key2",
										Value: MustPack("oldpayload"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			want:    []DiffItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diff(tt.args.src, tt.args.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%q. Diff():\n-got/+want\ndiff %s", tt.name, diff)
			}
		})
	}
}

func TestDiff_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		var (
			src bundlev1.Bundle
			dst bundlev1.Bundle
		)

		// Prepare arguments
		f.Fuzz(&src)
		f.Fuzz(&dst)

		// Execute
		Diff(&src, &dst)
	}
}
