// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"regexp"
)

// matchPathRule returns true if input match one of regexp.
func matchPathRule(input string, regexps []*regexp.Regexp) bool {
	match := false
	for _, r := range regexps {
		if r.MatchString(input) {
			match = true
			break
		}
	}

	return match
}

// collect generates an output array according to the strategy selected.
func collect(input []string, regexps []*regexp.Regexp, appendIfMatch bool) []string {
	out := []string{}

	for _, p := range input {
		if matchPathRule(p, regexps) == appendIfMatch {
			out = append(out, p)
		}
	}

	return out
}
