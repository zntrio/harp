// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/secret"
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

func Test_Bundle_DumpLoad(t *testing.T) {
	testCases := []struct {
		name    string
		input   *bundlev1.Bundle
		wantErr bool
	}{
		{
			name:    "Nil bundle",
			wantErr: true,
		},
		{
			name:    "Empty bundle",
			input:   &bundlev1.Bundle{},
			wantErr: false,
		},
		{
			name: "Filled bundle",
			input: &bundlev1.Bundle{
				Version: 1,
				Packages: []*bundlev1.Package{
					{
						Name: "infra/aws/foo/us-east-1/rds/postgresql/root_credentials",
						Secrets: &bundlev1.SecretChain{
							Version: 0,
							Data: []*bundlev1.KV{
								{
									Key:   "database_root_password",
									Type:  "string",
									Value: secret.MustPack("foo"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		testCase := tC
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()

			output := bytes.NewBuffer(nil)
			err := Dump(output, testCase.input)
			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Dump call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			inputTree, inputStats, err := Tree(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Tree call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			got, err := Load(output)
			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Load call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			outputTree, outputStats, err := Tree(got)
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Tree verification all, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			if !cmp.Equal(outputTree.Root(), inputTree.Root()) {
				t.Errorf("merkle tree root are different")
				return
			}

			if !cmp.Equal(outputStats.SecretCount, inputStats.SecretCount) {
				t.Errorf("secret count are different")
				return
			}

			if diff := cmp.Diff(got, testCase.input, ignoreOpts...); diff != "" {
				t.Errorf("%q. Bundle.Load():\n-got/+want\ndiff %s", testCase.name, diff)
			}
		})
	}
}

func Test_Bundle_JSONDumpLoad(t *testing.T) {
	testCases := []struct {
		name    string
		input   *bundlev1.Bundle
		wantErr bool
	}{
		{
			name:    "Nil bundle",
			wantErr: true,
		},
		{
			name:    "Empty bundle",
			input:   &bundlev1.Bundle{},
			wantErr: false,
		},
		{
			name: "Filled bundle",
			input: &bundlev1.Bundle{
				Version: 1,
				Packages: []*bundlev1.Package{
					{
						Name: "infra/aws/foo/us-east-1/rds/postgresql/root_credentials",
						Secrets: &bundlev1.SecretChain{
							Version: 0,
							Data: []*bundlev1.KV{
								{
									Key:   "database_root_password",
									Type:  "string",
									Value: secret.MustPack("foo"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		testCase := tC
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()

			output := bytes.NewBuffer(nil)
			err := AsProtoJSON(output, testCase.input)
			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the JSON call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			inputTree, inputStats, err := Tree(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Tree call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			got, err := FromDump(output)
			// Assert results expectations
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Load call, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			outputTree, outputStats, err := Tree(got)
			if (err != nil) != testCase.wantErr {
				t.Errorf("error during the Tree verification all, error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}

			if diff := cmp.Diff(testCase.input, got, ignoreOpts...); diff != "" {
				t.Errorf("%q. Bundle.FromDump():\n-got/+want\ndiff %s", testCase.name, diff)
			}

			if !cmp.Equal(outputStats.SecretCount, inputStats.SecretCount) {
				t.Errorf("secret count are different")
				return
			}

			if !cmp.Equal(outputTree.Root(), inputTree.Root()) {
				t.Errorf("merkle tree root are different")
				return
			}
		})
	}
}
