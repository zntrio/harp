// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hcl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/require"
)

func init() {
	goldie.FixtureDir = "testdata"
	spew.Config.DisablePointerAddresses = true
}

func TestParseFile(t *testing.T) {
	f, err := os.Open("testdata")
	require.NoError(t, err)
	defer f.Close()

	fis, err := f.Readdir(-1)
	require.NoError(t, err)
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if filepath.Ext(fi.Name()) == ".golden" {
			continue
		}

		t.Run(fi.Name(), func(t *testing.T) {
			_, err := ParseFile(filepath.Join("testdata", fi.Name()))
			require.NoError(t, err)

			// goldie.Assert(t, fi.Name(), []byte(spew.Sdump(cfg)))
		})
	}
}
