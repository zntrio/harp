// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.18
// +build go1.18

package loader_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"zntr.io/harp/v2/pkg/bundle/template"
)

func loadFromFile(t testing.TB, filename string) []byte {
	t.Helper()

	// Load sample
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unable to load content '%v': %v", filename, err)
	}

	// Load all content
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("unable to load all content '%v': %v", filename, err)
	}

	return content
}

func FuzzBundleLoader(f *testing.F) {
	f.Add(loadFromFile(f, "../../fixtures/template/valid/blank.yaml"))
	f.Add(loadFromFile(f, "../../../samples/customer-bundle/spec.yaml"))

	f.Fuzz(func(t *testing.T, in []byte) {
		// Read from randomized data
		template.YAML(bytes.NewBuffer(in))
	})
}
