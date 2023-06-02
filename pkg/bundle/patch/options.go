// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package patch

type options struct {
	stopAtRuleID      string
	stopAtRuleIndex   int
	ignoreRuleIDs     []string
	ignoreRuleIndexes []int
}

type OptionFunc func(o *options)

// -----------------------------------------------------------------------------

// WithStopAtRuleID sets the id of the rule to stop evaluation at.
func WithStopAtRuleID(value string) OptionFunc {
	return func(o *options) {
		o.stopAtRuleID = value
	}
}

// WithStopAtRuleIndex sets the index of the rule to stop evaluation at.
func WithStopAtRuleIndex(value int) OptionFunc {
	return func(o *options) {
		o.stopAtRuleIndex = value
	}
}

// WithIgnoreRuleIDs sets the rule identifiers to ignore.
func WithIgnoreRuleIDs(values ...string) OptionFunc {
	return func(o *options) {
		o.ignoreRuleIDs = values
	}
}

// WithIgnoreRuleIndexes sets the rule indexes to ignore.
func WithIgnoreRuleIndexes(values ...int) OptionFunc {
	return func(o *options) {
		o.ignoreRuleIndexes = values
	}
}
