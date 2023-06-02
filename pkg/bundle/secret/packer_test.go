// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secret

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
)

func Test_Pack_Pack_Unpack(t *testing.T) {
	testCases := []struct {
		desc    string
		in      interface{}
		wantErr bool
	}{
		{
			desc:    "empty struct",
			in:      map[interface{}]interface{}{},
			wantErr: true,
		},
		{
			desc:    "string",
			in:      "foo",
			wantErr: false,
		},
		{
			desc:    "bytes",
			in:      []byte("foo"),
			wantErr: false,
		},
		{
			desc:    "uint8 array",
			in:      []uint8("foo"),
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := Pack(tC.in)
			// Assert results expectations
			if (err != nil) != tC.wantErr {
				t.Errorf("error during the call, error = %v, wantErr %v", err, tC.wantErr)
				return
			}

			var out interface{}
			err = Unpack(got, &out)
			// Assert results expectations
			if (err != nil) != tC.wantErr {
				t.Errorf("error during the call, error = %v, wantErr %v", err, tC.wantErr)
				return
			}

			if tC.wantErr {
				return
			}

			if diff := cmp.Diff(out, tC.in); diff != "" {
				t.Errorf("%q. Secret.PackUnpack():\n-got/+want\ndiff %s", tC.desc, diff)
			}
		})
	}
}

func Test_UnPack_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 100000; i++ {
		f := fuzz.New()

		var (
			in  []byte
			out struct{}
		)

		// Fuzz input
		f.Fuzz(&in)

		// Execute
		Unpack(in, &out)
	}
}
