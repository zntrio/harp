// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package encryption

import (
	"errors"
	"fmt"
	"strings"

	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/sdk/value"
)

// FromKey returns the value transformer that match the value format.
func FromKey(keyValue string) (value.Transformer, error) {
	var (
		transformer value.Transformer
		err         error
	)

	// Check arguments
	if keyValue == "" {
		return nil, fmt.Errorf("unable to select a value transformer with blank value")
	}

	// Extract prefix
	parts := strings.SplitN(keyValue, ":", 2)
	if len(parts) != 2 {
		// Fallback to fernet
		parts = []string{"fernet", keyValue}
	}

	// Clean prefix
	prefix := strings.ToLower(strings.TrimSpace(parts[0]))

	// Build the value transformer according to used prefix.
	tf, ok := registry[prefix]
	if !ok {
		return nil, fmt.Errorf("no transformer registered for %q as prefix", prefix)
	}

	// Build the transformer instance
	transformer, err = tf(keyValue)

	// Check transformer initialization error
	if transformer == nil || err != nil {
		return nil, fmt.Errorf("unable to initialize value transformer: %w", err)
	}

	// No error
	return transformer, nil
}

// Must is used to panic when a transformer initialization failed.
func Must(t value.Transformer, err error) value.Transformer {
	if err != nil {
		panic(err)
	}
	if types.IsNil(t) {
		panic(errors.New("transformer is nil with a nil error"))
	}

	return t
}
