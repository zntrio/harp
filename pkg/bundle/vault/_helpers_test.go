// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"regexp"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func Test_matchPathRule_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var (
			value   string
			regexps []*regexp.Regexp
		)

		f.Fuzz(&value)
		f.Fuzz(&regexps)

		// Execute
		matchPathRule(value, regexps)
	}
}

func Test_collect_Fuzz(t *testing.T) {
	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var (
			values        []string
			regexps       []*regexp.Regexp
			appendIfMatch bool
		)

		f.Fuzz(&values)
		f.Fuzz(&regexps)
		f.Fuzz(&appendIfMatch)

		// Execute
		collect(values, regexps, appendIfMatch)
	}
}
