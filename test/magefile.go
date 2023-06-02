// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
