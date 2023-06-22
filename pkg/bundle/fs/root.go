// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build go1.16
// +build go1.16

package fs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/awnumar/memguard"
	"google.golang.org/protobuf/proto"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

const (
	directoryAccess = 0o555
	fileAccess      = 0o444
)

type bundleFs struct {
	root *directory
}

// -----------------------------------------------------------------------------

// FromBundle initializes an fs.FS object from the given bundle.
func FromBundle(b *bundlev1.Bundle) (BundleFS, error) {
	// Check arguments
	if b == nil {
		return nil, errors.New("unable to create a filesytem from a nil bundle")
	}

	// Prepare vfs root
	bfs := &bundleFs{
		root: &directory{
			children: map[string]interface{}{},
		},
	}

	// Prepare filesystem
	for _, p := range b.Packages {
		if p == nil {
			// ignore nil package
			continue
		}

		// Serialize package
		body, err := proto.Marshal(p)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize package %q: %w", p.Name, err)
		}

		// Write content
		if errWrite := bfs.WriteFile(p.Name, body, fileAccess); errWrite != nil {
			return nil, fmt.Errorf("unable to write package %q in filesystem: %w", p.Name, err)
		}
	}

	// Return bundle filesystem
	return bfs, nil
}

// -----------------------------------------------------------------------------

func (bfs *bundleFs) Open(name string) (fs.File, error) {
	// Return root as default
	if name == "" {
		return bfs.root, nil
	}

	// Validate input path
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}

	// Create directory tree
	dirPath, name := filepath.Split(name)
	dirNames := strings.Split(dirPath, "/")

	// Browse directory tree
	currentDirectory := bfs.root
	for _, dirName := range dirNames {
		// Skip empty directory name
		if dirName == "" {
			continue
		}

		it, ok := currentDirectory.children[dirName]
		if !ok {
			return nil, fmt.Errorf("directory %q not found: %w", dirName, fs.ErrNotExist)
		}
		currentDirectory, ok = it.(*directory)
		if !ok {
			return nil, errors.New("invalid directory iterator value")
		}
	}

	// Get child
	h, ok := currentDirectory.children[name]
	if !ok {
		return nil, fmt.Errorf("item %q not found in directory %q: %w", name, currentDirectory.name, fs.ErrNotExist)
	}

	switch it := h.(type) {
	case *directory:
		// Return directory
		return it, nil
	case *file:
		// Open enclave
		body, err := it.content.Open()
		if err != nil {
			return nil, fmt.Errorf("file %q could not be opened: %w", name, err)
		}

		// Assign body reader
		it.bodyReader = body.Reader()

		// Return file
		return it, nil
	}

	return nil, fmt.Errorf("unexpected file type in filesystem %s: %w", name, fs.ErrInvalid)
}

func (bfs *bundleFs) ReadDir(name string) ([]fs.DirEntry, error) {
	// Try to open directory
	h, err := bfs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open directory: %w", err)
	}

	// Retrieve directory info
	fi, err := h.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve directory info: %w", err)
	}

	// Confirm it's a directory
	if !fi.IsDir() {
		return nil, fmt.Errorf("path %q point to a file", name)
	}

	// Convert handle to directory reader
	dir, ok := h.(fs.ReadDirFile)
	if !ok {
		return nil, fmt.Errorf("path %q point to a directory but could not be listed", name)
	}

	// Delegate to directory list
	return dir.ReadDir(0)
}

func (bfs *bundleFs) ReadFile(name string) ([]byte, error) {
	// Try to open file
	h, err := bfs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	// Delete to file
	return io.ReadAll(h)
}

func (bfs *bundleFs) WriteFile(name string, data []byte, perm os.FileMode) error {
	// Create directory tree
	dirPath, fname := filepath.Split(name)
	dirNames := strings.Split(dirPath, "/")

	// MkDirAll
	currentDirectory := bfs.root
	for _, dirName := range dirNames {
		// Skip empty directory name
		if dirName == "" {
			continue
		}

		currentDirectory.RLock()
		it, ok := currentDirectory.children[dirName]
		currentDirectory.RUnlock()
		if !ok {
			it = &directory{
				name:     dirName,
				perm:     directoryAccess,
				children: map[string]interface{}{},
			}

			currentDirectory.Lock()
			currentDirectory.children[dirName] = it
			currentDirectory.Unlock()
		}
		currentDirectory, ok = it.(*directory)
		if !ok {
			return errors.New("invalid directory iterator value")
		}
	}

	// Create file entry
	currentDirectory.Lock()
	currentDirectory.children[fname] = &file{
		name:    fname,
		mode:    fileAccess,
		size:    int64(len(data)),
		content: memguard.NewEnclave(data),
	}
	currentDirectory.Unlock()

	// No error
	return nil
}

func (bfs *bundleFs) Stat(name string) (fs.FileInfo, error) {
	// Try to open file
	h, err := bfs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	// Delegate to file
	return h.Stat()
}
