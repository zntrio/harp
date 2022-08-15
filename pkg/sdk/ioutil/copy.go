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

package ioutil

import (
	"errors"
	"io"
)

var (
	// ErrTruncatedCopy is raised when the copy is larger than expected.
	ErrTruncatedCopy = errors.New("truncated copy due to too large input")
)

// Copy uses a buffered CopyN and a hardlimit to stop read from the reader when
// the maxSize amount of data has been writtent to the given writer.
func Copy(maxSize int64, w io.Writer, r io.Reader) error {
	var (
		contentLength = int64(0)
	)

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
