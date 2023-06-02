// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

import (
	"fmt"
)

// SecretReaderFunc is a function to retrieve a secret from a given path.
type SecretReaderFunc func(path string) (map[string]interface{}, error)

// SecretReaders uses given secret reader funcs to resolve secret path.
func SecretReaders(secretReaders []SecretReaderFunc) func(string) (map[string]interface{}, error) {
	return func(secretPath string) (map[string]interface{}, error) {
		// For all secret readers
		for _, sr := range secretReaders {
			value, err := sr(secretPath)
			if err != nil {
				// Check next secret reader
				continue
			}

			// No error
			return value, nil
		}

		// Return error
		return nil, fmt.Errorf("no value found for %q, check secret path or secret reader settings", secretPath)
	}
}
