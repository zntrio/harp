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

	"zntr.io/harp/v2/pkg/bundle"
)

func loadFromFile(t testing.TB, filename string) []byte {
	t.Helper()

	// Load sample bundle
	completeBundle, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unable to load bundle content '%v': %v", filename, err)
	}

	// Load all content
	content, err := io.ReadAll(completeBundle)
	if err != nil {
		t.Fatalf("unable to load all bundle content '%v': %v", filename, err)
	}

	return content
}

func FuzzBundleLoader(f *testing.F) {
	f.Add(loadFromFile(f, "../../../fixtures/bundles/complete.bundle"))
	f.Add(loadFromFile(f, "../../../fixtures/bundles/empty.bundle"))

	f.Fuzz(func(t *testing.T, in []byte) {
		// Read from randomized data
		bundle.FromContainerReader(bytes.NewBuffer(in))
	})
}
