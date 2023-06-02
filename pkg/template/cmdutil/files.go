// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmdutil

import (
	"fmt"
	"io/fs"

	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/template/engine"
	"zntr.io/harp/v2/pkg/template/files"
)

// Files returns template files.
func Files(fileSystem fs.FS, basePath string) (engine.Files, error) {
	// Check arguments
	if types.IsNil(fileSystem) {
		return nil, fmt.Errorf("unable to load files without a default filesystem to use")
	}

	// Get appropriate loader
	loader, err := files.Loader(fileSystem, basePath)
	if err != nil {
		return nil, fmt.Errorf("unable to process files: %w", err)
	}

	// Crawl and load file content
	fileList, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}

	// Wrap as template files
	templateFiles := engine.NewFiles(fileList)

	// No error
	return templateFiles, nil
}
