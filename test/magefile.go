// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"

	"zntr.io/harp/v2/build/mage/golang"
)

// -----------------------------------------------------------------------------

type Fuzz mg.Namespace

func (Fuzz) BundleLoader() {
	mg.SerialDeps(
		golang.FuzzBuild("bundle-loader", "zntr.io/harp/v2/test/fuzz/bundle/loader"),
		golang.FuzzRun("bundle-loader"),
	)
}

func (Fuzz) TemplateReader() {
	mg.SerialDeps(
		golang.FuzzBuild("template-reader", "zntr.io/harp/v2/test/fuzz/template"),
		golang.FuzzRun("template-reader"),
	)
}
