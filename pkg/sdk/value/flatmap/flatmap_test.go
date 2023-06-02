// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package flatmap

import (
	"reflect"
	"testing"

	"zntr.io/harp/v2/pkg/bundle"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		Input  map[string]interface{}
		Output map[string]bundle.KV
	}{
		{
			Input: map[string]interface{}{
				"app": bundle.KV{
					"database": bundle.KV{
						"user":     "app-12345679",
						"password": "testpassword",
					},
					"integers": bundle.KV{
						"int":    int(8),
						"int8":   int8(0x7F),
						"int16":  int16(0x7FF),
						"int32":  int32(0x7FFFFFFF),
						"int64":  int64(0x7FFFFFFFFFFFFFFF),
						"uint":   uint(8),
						"uint8":  uint8(0xFF),
						"uint16": uint16(0xFFF),
						"uint32": uint32(0xFFFFFFFF),
						"uint64": uint64(0xFFFFFFFFFFFFFFFF),
					},
					"arrays": bundle.KV{
						"strings": []string{"first", "second", "third"},
					},
				},
			},
			Output: map[string]bundle.KV{
				"app/database": {
					"user":     "app-12345679",
					"password": "testpassword",
				},
				"app/integers": {
					"int":    "8",
					"int16":  "2047",
					"int32":  "2147483647",
					"int64":  "9223372036854775807",
					"int8":   "127",
					"uint":   "8",
					"uint16": "4095",
					"uint32": "4294967295",
					"uint64": "18446744073709551615",
					"uint8":  "255",
				},
				"app/arrays/strings": {
					"#": "3",
					"0": "first",
					"1": "second",
					"2": "third",
				},
			},
		},
	}

	for _, tc := range cases {
		actual := Flatten(tc.Input)
		if !reflect.DeepEqual(actual, tc.Output) {
			t.Fatalf(
				"Input:\n\n%#v\n\nOutput:\n\n%#v\n\nExpected:\n\n%#v\n",
				tc.Input,
				actual,
				tc.Output)
		}
	}
}
