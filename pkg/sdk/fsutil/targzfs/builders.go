// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package targzfs

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/zntrio/harp/v2/pkg/sdk/ioutil"
)

// FromFile creates an archive filesystem from a filename.
func FromFile(name string) (fs.FS, error) {
	// Open the target file
	fn, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open archive %q: %w", name, err)
	}

	// Delegate to reader constructor.
	return FromReader(fn)
}

// FromReader exposes the contents of the given reader (which is a .tar.gz file)
// as an fs.FS.
func FromReader(r io.Reader) (fs.FS, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("unable to open .tar.gz file: %w", err)
	}

	// Retrieve TAR content from GZIP
	var (
		tarContents bytes.Buffer
	)

	// Chunked read with hard limit to prevent/reduce zipbomb vulnerability
	// exploitation.
	if err := ioutil.Copy(maxDecompressedSize, &tarContents, gz); err != nil {
		return nil, fmt.Errorf("unable to decompress the archive: %w", err)
	}

	// Close the gzip decompressor
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("unable to close gzip reader: %w", err)
	}

	// TAR format reader
	tarReader := tar.NewReader(&tarContents)

	// Prepare in-memory filesystem.
	ret := &tarGzFs{
		files:       make(map[string]*tarEntry),
		rootEntries: make([]fs.DirEntry, 0, 10),
		rootEntry:   nil,
	}

	for {
		// Iterate on each file entry
		hdr, err := tarReader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("unable to read .tar.gz entry: %w", err)
		}
		if hdr != nil && len(ret.files) > maxFileCount {
			return nil, errors.New("interrupted extraction, too many files in the archive")
		}

		// Clean file path. (ZipSlip)
		name := path.Clean(hdr.Name)
		if name == "." {
			continue
		}

		// Load content in memory
		var (
			fileContents bytes.Buffer
		)

		// Chunked read with hard limit to prevent/reduce post decompression
		// explosion
		if err := ioutil.Copy(maxFileSize, &fileContents, tarReader); err != nil {
			return nil, fmt.Errorf("unable to copy file content to memory: %w", err)
		}

		// Register file
		e := &tarEntry{
			h:       hdr,
			b:       fileContents.Bytes(),
			entries: nil,
		}

		// Add as file entry
		ret.files[name] = e

		// Create directories
		dir := path.Dir(name)
		if dir == "." {
			ret.rootEntries = append(ret.rootEntries, e)
		} else {
			if parent, ok := ret.files[dir]; ok {
				parent.entries = append(parent.entries, e)
			}
		}
	}

	// No error
	return ret, nil
}
