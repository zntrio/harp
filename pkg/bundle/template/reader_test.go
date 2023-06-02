// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package template

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func mustLoad(filePath string) io.Reader {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return f
}

type readerTestCase struct {
	name    string
	args    io.Reader
	wantErr bool
}

//nolint:unparam // rootPath has always the same value for the moment
func generateReaderTests(t *testing.T, rootPath, state string, wantErr bool) []readerTestCase {
	tests := []readerTestCase{}
	// Generate invalid test cases
	if err := filepath.Walk(filepath.Join(rootPath, state), func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}
		if info.IsDir() {
			return nil
		}

		tests = append(tests, readerTestCase{
			name:    fmt.Sprintf("%s-%s", state, filepath.Base(info.Name())),
			args:    mustLoad(path),
			wantErr: wantErr,
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	return tests
}

func TestYAML(t *testing.T) {
	tests := []readerTestCase{
		{
			name:    "nil",
			wantErr: true,
		},
	}

	// Generate invalid test cases
	tests = append(tests, generateReaderTests(t, "../../../test/fixtures/template", "invalid", true)...)

	// Generate valid test cases
	tests = append(tests, generateReaderTests(t, "../../../test/fixtures/template", "valid", false)...)

	// Execute them
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := YAML(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("YAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
