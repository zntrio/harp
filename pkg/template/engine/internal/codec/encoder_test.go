// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package codec

import (
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestToYAML(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				v: map[string]string{},
			},
			want: "{}",
		},
		{
			name: "object",
			args: args{
				v: map[string]string{
					"key": "value",
				},
			},
			want: "key: value",
		},
		{
			name: "non-serializable",
			args: args{
				v: map[string]interface{}{
					"key": make(chan string, 1),
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToYAML(tt.args.v); got != tt.want {
				t.Errorf("ToYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTOML(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				v: map[string]string{},
			},
			want: "",
		},
		{
			name: "object",
			args: args{
				v: map[string]string{
					"key": "value",
				},
			},
			want: "key = \"value\"\n",
		},
		/*{
			name: "non-serializable",
			args: args{
				v: map[string]interface{}{
					"key": make(chan string, 1),
				},
			},
			want: "",
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTOML(tt.args.v); got != tt.want {
				t.Errorf("ToTOML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				v: map[string]string{},
			},
			want: "{}",
		},
		{
			name: "object",
			args: args{
				v: map[string]string{
					"key": "value",
				},
			},
			want: `{"key":"value"}`,
		},
		{
			name: "non-serializable",
			args: args{
				v: map[string]interface{}{
					"key": make(chan string, 1),
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.args.v); got != tt.want {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

// -----------------------------------------------------------------------------

func TestToYAML_Fuzz(t *testing.T) {
	// Making sure that it never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var input struct {
			Integer int
			String  string
			Map     map[string]string
		}

		// Fuzz input
		f.Fuzz(&input)

		// Execute
		ToYAML(input)
	}
}

func TestToTOML_Fuzz(t *testing.T) {
	// Making sure that it never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var input struct {
			Integer int
			String  string
			Map     map[string]string
		}

		// Fuzz input
		f.Fuzz(&input)

		// Execute
		ToTOML(input)
	}
}

func TestToJSON_Fuzz(t *testing.T) {
	// Making sure that it never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		var input struct {
			Integer int
			String  string
			Map     map[string]string
		}

		// Fuzz input
		f.Fuzz(&input)

		// Execute
		ToJSON(input)
	}
}
