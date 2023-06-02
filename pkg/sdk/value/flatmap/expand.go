// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package flatmap

import (
	"strings"

	"zntr.io/harp/v2/pkg/bundle"
)

// Expand takes a map and a key (prefix) and expands that value into
// a more complex structure. This is the reverse of the Flatten operation.
func Expand(m bundle.KV, key string) interface{} {
	// If the key is exactly a key in the map, just return it
	if v, ok := m[key]; ok {
		if v == "true" {
			return true
		} else if v == "false" {
			return false
		}

		return v
	}

	// Check if this is a prefix in the map
	prefix := key
	if key != "" {
		prefix = key + "/"
	}
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			return expandMap(m, prefix)
		}
	}

	return nil
}

func expandMap(m bundle.KV, prefix string) bundle.KV {
	result := make(bundle.KV)
	for k := range m {
		if !strings.HasPrefix(k, prefix) {
			// Prefix not found
			continue
		}

		// Remove the prefix
		key := k[len(prefix):]
		idx := strings.Index(key, "/")
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}

		// Recursive call to handle subtree
		result[key] = Expand(m, k[:len(prefix)+len(key)])
	}

	return result
}
