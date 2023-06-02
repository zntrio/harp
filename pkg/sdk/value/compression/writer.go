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

// NewWriter returns a wrtier implementation according to given algorithm.
func NewWriter(w io.Writer, algorithm string) (io.WriteCloser, error) {
	// Normalize input
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	var (
		compressedWriter io.WriteCloser
		writerErr        error
	)

	// Apply transformation
	switch algorithm {
	case "gzip":
		compressedWriter = gzip.NewWriter(w)
	case "lzw", "lzw-lsb":
		compressedWriter = lzw.NewWriter(w, lzw.LSB, 8)
	case "lzw-msb":
		compressedWriter = lzw.NewWriter(w, lzw.MSB, 8)
	case "lz4":
		compressedWriter = lz4.NewWriter(w)
	case "s2", "snappy":
		compressedWriter = s2.NewWriter(w)
	case "zlib":
		compressedWriter = zlib.NewWriter(w)
	case "flate":
		compressedWriter, writerErr = flate.NewWriter(w, flate.DefaultCompression)
		if writerErr != nil {
			return nil, fmt.Errorf("unable to initialize flate compressor: %w", writerErr)
		}
	case "lzma":
		compressedWriter, writerErr = xz.NewWriter(w)
		if writerErr != nil {
			return nil, fmt.Errorf("unable to initialize lzma compressor: %w", writerErr)
		}
	case "zstd":
		compressedWriter, writerErr = zstd.NewWriter(w)
		if writerErr != nil {
			return nil, fmt.Errorf("unable to initialize zstd compressor: %w", writerErr)
		}
	default:
		return nil, fmt.Errorf("unhandled compression algorithm %q", algorithm)
	}

	// No error
	return compressedWriter, nil
}
