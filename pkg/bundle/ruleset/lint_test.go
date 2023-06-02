// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ruleset

import "testing"

func TestLint(t *testing.T) {
	tests := []readerTestCase{
		{
			name:    "nil",
			wantErr: true,
		},
	}

	// Generate invalid test cases
	tests = append(tests, generateReaderTests(t, "../../../test/fixtures/ruleset", "invalid", true)...)

	// Generate valid test cases
	tests = append(tests, generateReaderTests(t, "../../../test/fixtures/ruleset", "valid", false)...)

	// Execute them
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Lint(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
