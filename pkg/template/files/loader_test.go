// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDir(t *testing.T) {
	// Get basepath
	basePath, err := filepath.Abs("../../../test")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}

	// Initialize filesystem
	fileSystem := os.DirFS(basePath)

	l, err := Loader(fileSystem, ".")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	c, err := l.Load()
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	if len(c) == 0 {
		t.Fatalf("Failed to load all test files")
	}
}
