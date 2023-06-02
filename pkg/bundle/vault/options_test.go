// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestWithExcludePath_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var value string
		f.Fuzz(&value)

		// Execute
		opts := &options{}
		WithExcludePath(value)(opts)
	}
}

func TestWithIncludePath_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var value string
		f.Fuzz(&value)

		// Execute
		opts := &options{}
		WithIncludePath(value)(opts)
	}
}

func TestWithPrefix_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var value string
		f.Fuzz(&value)

		// Execute
		opts := &options{}
		WithPrefix(value)(opts)
	}
}
