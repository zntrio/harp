// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package fsutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/log"
)

// Dump the given vfs to the outputpath.
func Dump(srcFs fs.FS, outPath string) error {
	return fs.WalkDir(srcFs, ".", func(path string, d fs.DirEntry, errWalk error) error {
		// Raise immediately the error if any.
		if errWalk != nil {
			return fmt.Errorf("%s: %w", path, errWalk)
		}

		// Ignore directory
		if d.IsDir() {
			return nil
		}

		// Compute the target path
		targetPath := filepath.Join(outPath, path)

		// Extract relative directory
		relativeDir := filepath.Dir(targetPath)

		// Check folder hierarchy existence.
		if _, err := os.Stat(relativeDir); os.IsNotExist(err) {
			if err := os.MkdirAll(relativeDir, 0o750); err != nil {
				return fmt.Errorf("unable to create intermediate directories for path %q: %w", relativeDir, err)
			}
		}

		// Encsure not out of safe directory file creation.
		cleanTargetPath := filepath.Clean(targetPath)
		if !strings.HasPrefix(cleanTargetPath, outPath) {
			return fmt.Errorf("unable to create %q file, the path is not in the expected output path", targetPath)
		}

		// Create file
		targetFile, err := os.Create(cleanTargetPath)
		if err != nil {
			return fmt.Errorf("unable to create the output file: %w", err)
		}

		// Open input file
		srcFile, err := srcFs.Open(path)
		if err != nil {
			return fmt.Errorf("unable to open source file: %w", err)
		}

		log.Bg().Debug("Copy file ...", zap.String("file", path))

		// Open the target file
		if _, err := io.Copy(targetFile, srcFile); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("unable to copy content from %q to %q: %w", path, targetPath, err)
			}
		}

		// Close the file
		return srcFile.Close()
	})
}
