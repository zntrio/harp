// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"errors"
)

// ErrKeyNotFound is raised when the given key could not be found in the store.
var ErrKeyNotFound = errors.New("key not found")

// Store describes the key/value store contract.
type Store interface {
	// Get the value stored at the given key.
	Get(ctx context.Context, key string) (*Pair, error)
	// Exists checks if the key exists inside the store
	Exists(ctx context.Context, key string) (bool, error)
	// Delete a value addressed by "key"
	Delete(ctx context.Context, key string) error
	// Put the given value at the given key.
	Put(ctx context.Context, key string, value []byte) error
	// List subkeys at a given path
	List(ctx context.Context, path string) ([]*Pair, error)
	// Close closes the client connection
	Close() error
}

// -----------------------------------------------------------------------------

type Pair struct {
	Key     string
	Value   []byte
	Version uint64
}
