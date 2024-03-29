// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package diceware

import (
	"strings"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestDiceware(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantCount int
	}{
		{
			name: "negative",
			args: args{
				count: -1,
			},
			wantErr:   false,
			wantCount: MinWordCount,
		},
		{
			name: "zero",
			args: args{
				count: 0,
			},
			wantErr:   false,
			wantCount: MinWordCount,
		},
		{
			name: "five",
			args: args{
				count: 5,
			},
			wantErr:   false,
			wantCount: 5,
		},
		{
			name: "upper limit",
			args: args{
				count: MaxWordCount + 1,
			},
			wantErr:   false,
			wantCount: MaxWordCount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diceware(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diceware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotWordCount := len(strings.Split(got, "-"))
			if (tt.wantCount > 0) && tt.wantCount != gotWordCount {
				t.Errorf("Diceware() expected word count = %v, got %v", tt.wantCount, gotWordCount)
				return
			}
		})
	}
}

func TestPredefined(t *testing.T) {
	tests := []struct {
		name      string
		callable  func() (string, error)
		wantCount int
		wantErr   bool
	}{
		{
			name:      "basic",
			callable:  Basic,
			wantCount: BasicWordCount,
			wantErr:   false,
		},
		{
			name:      "strong",
			callable:  Strong,
			wantCount: StrongWordCount,
			wantErr:   false,
		},
		{
			name:      "paranoid",
			callable:  Paranoid,
			wantCount: ParanoidWordCount,
			wantErr:   false,
		},
		{
			name:      "master",
			callable:  Master,
			wantCount: MasterWordCount,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.callable()
			if (err != nil) != tt.wantErr {
				t.Errorf("Predefined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotWordCount := len(strings.Split(got, "-"))
			if (tt.wantCount > 0) && tt.wantCount != gotWordCount {
				t.Errorf("Predefined() expected word count = %v, got %v", tt.wantCount, gotWordCount)
				return
			}
		})
	}
}

// -----------------------------------------------------------------------------

func TestDiceware_Fuzz(t *testing.T) {
	// Making sure that it never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var wordCount int

		// Fuzz input
		f.Fuzz(&wordCount)

		// Execute
		Diceware(wordCount)
	}
}
