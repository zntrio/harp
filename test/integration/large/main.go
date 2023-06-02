// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package main

import (
	"fmt"
	"os"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/secret"
)

func main() {
	b := &bundlev1.Bundle{
		Packages: []*bundlev1.Package{},
	}

	// Create 25000 packages
	for i := 0; i < 25000; i++ {
		p := &bundlev1.Package{
			Name: fmt.Sprintf("app/secret/large-bundle/%d", i),
			Secrets: &bundlev1.SecretChain{
				Data: []*bundlev1.KV{},
			},
		}

		for j := 0; j < 100; j++ {
			p.Secrets.Data = append(p.Secrets.Data, &bundlev1.KV{
				Key:   fmt.Sprintf("secret-%d", j),
				Value: secret.MustPack("test-value"),
			})
		}

		b.Packages = append(b.Packages, p)
	}

	// Save as a container in Stdout.
	if err := bundle.ToContainerWriter(os.Stdout, b); err != nil {
		panic(err)
	}
}
