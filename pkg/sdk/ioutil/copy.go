// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ioutil

import (
	"errors"
	"io"
)

// ErrTruncatedCopy is raised when the copy is larger than expected.
var ErrTruncatedCopy = errors.New("truncated copy due to too large input")

// Copy uses a buffered CopyN and a hardlimit to stop read from the reader when
// the maxSize amount of data has been written to the given writer.
func Copy(maxSize int64, w io.Writer, r io.Reader) error {
	contentLength := int64(0)

	// Chunked read with hard limit to prevent/reduce zipbomb vulnerability
	// exploitation.
	for {
		written, err := io.CopyN(w, r, 1024)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		// Add to length
		contentLength += written

		// Check max size
		if contentLength > maxSize {
			return ErrTruncatedCopy
		}
	}

	// No error
	return nil
}
