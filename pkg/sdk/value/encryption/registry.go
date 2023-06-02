// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package encryption

import (
	"fmt"

	"zntr.io/harp/v2/pkg/sdk/value"
)

// TransformerFactoryFunc is used for transformer building for encryption.
type TransformerFactoryFunc func(string) (value.Transformer, error)

var registry map[string]TransformerFactoryFunc

// Register a transformer with the given prefix.
func Register(prefix string, factory TransformerFactoryFunc) {
	// Lazy initialization
	if registry == nil {
		registry = map[string]TransformerFactoryFunc{}
	}

	// Check if not already registered
	if _, ok := registry[prefix]; ok {
		panic(fmt.Errorf("encryption transformer already registered fro %q prefix", prefix))
	}

	// Register the transformer
	registry[prefix] = factory
}
