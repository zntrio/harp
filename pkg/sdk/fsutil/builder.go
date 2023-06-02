// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package fsutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"zntr.io/harp/v2/pkg/sdk/fsutil/targzfs"
)

var ErrNotSupported = errors.New("not supported filesystem path")

func From(rootPath string) (fs.FS, error) {
	// Get absolute path
	absPath, errPath := filepath.Abs(filepath.Clean(rootPath))
	if errPath != nil {
		return nil, fmt.Errorf("unable to get absolute path: %w", errPath)
	}

	// Check existence
	fi, errFileInfo := os.Stat(absPath)
	if errFileInfo != nil {
		return nil, fmt.Errorf("unabe to retrieve file information: %w", errFileInfo)
	}

	var (
		errRootFS  error
		fileRootFS fs.FS
	)
	switch {
	case !fi.IsDir() && strings.HasSuffix(absPath, ".tar.gz"):
		// Create in-memory filesystem.
		fileRootFS, errRootFS = targzfs.FromFile(os.DirFS(filepath.Dir(absPath)), filepath.Base(absPath))
	case fi.IsDir():
		fileRootFS = os.DirFS(absPath)
	default:
		errRootFS = ErrNotSupported
	}
	if errRootFS != nil {
		return nil, fmt.Errorf("unable to prepare filesystem: %w", errRootFS)
	}

	// No error
	return fileRootFS, nil
}
