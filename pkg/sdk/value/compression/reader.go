// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package compression

import (
	"compress/lzw"
	"fmt"
	"io"
	"strings"

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"github.com/ulikunitz/xz"
)

// -----------------------------------------------------------------------------

// NewReader returns a writer implementation according to given algorithm.
func NewReader(r io.Reader, algorithm string) (io.ReadCloser, error) {
	// Normalize input
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	var (
		compressedReader io.ReadCloser
		readerErr        error
	)

	// Apply transformation
	switch algorithm {
	case "identity":
		compressedReader = io.NopCloser(r)
	case "gzip":
		compressedReader, readerErr = gzip.NewReader(r)
	case "lzw", "lzw-lsb":
		compressedReader = lzw.NewReader(r, lzw.LSB, 8)
	case "lzw-msb":
		compressedReader = lzw.NewReader(r, lzw.MSB, 8)
	case "lz4":
		compressedReader = io.NopCloser(lz4.NewReader(r))
	case "s2", "snappy":
		compressedReader = io.NopCloser(s2.NewReader(r))
	case "zlib":
		compressedReader, readerErr = zlib.NewReader(r)
	case "flate", "deflate":
		compressedReader = flate.NewReader(r)
	case "lzma":
		reader, err := xz.NewReader(r)
		if err != nil {
			readerErr = err
		} else {
			compressedReader = io.NopCloser(reader)
		}
	case "zstd":
		reader, err := zstd.NewReader(r)
		if err != nil {
			readerErr = err
		} else {
			compressedReader = reader.IOReadCloser()
		}
	default:
		return nil, fmt.Errorf("unhandled compression algorithm %q", algorithm)
	}
	if readerErr != nil {
		return nil, fmt.Errorf("unable to initialize %q compressor: %w", algorithm, readerErr)
	}

	// No error
	return compressedReader, nil
}
