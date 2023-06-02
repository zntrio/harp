// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hash

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

func NewMultiHash(r io.Reader, algorithms ...string) (map[string]string, error) {
	hashers := map[string]hash.Hash{}

	// Instantiate hashers
	for _, algo := range algorithms {
		// Create an hasher instance.
		h, err := NewHasher(algo)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize %q algorithm: %w", algo, err)
		}

		// Assign to hashers.
		hashers[algo] = h
	}

	// Copy to all hashers
	_, err := io.Copy(io.MultiWriter(hashToMultiWriter(hashers)), r)
	if err != nil {
		return nil, err
	}

	// Finalize
	res := make(map[string]string)
	for algo, v := range hashers {
		res[algo] = hex.EncodeToString(v.Sum(nil))
	}

	// No error
	return res, nil
}

// -----------------------------------------------------------------------------

func hashToMultiWriter(hashers map[string]hash.Hash) io.Writer {
	w := make([]io.Writer, 0, len(hashers))
	for _, v := range hashers {
		w = append(w, v)
	}
	return io.MultiWriter(w...)
}
