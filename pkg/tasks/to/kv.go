// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package to

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/kv"
	"zntr.io/harp/v2/pkg/tasks"
)

type PublishKVTask struct {
	_               struct{}
	ContainerReader tasks.ReaderProvider
	Store           kv.Store
	SecretAsKey     bool
	Prefix          string
}

func (t *PublishKVTask) Run(ctx context.Context) error {
	// Create the reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input bundle reader: %w", err)
	}

	// Extract bundle from container
	b, err := bundle.FromContainerReader(reader)
	if err != nil {
		return fmt.Errorf("unable to load bundle: %w", err)
	}

	// Convert as map
	bundleMap, err := bundle.AsMap(b)
	if err != nil {
		return fmt.Errorf("unable to transform the bundle as a map: %w", err)
	}

	// Foreach element in the bundle map.
	for key, value := range bundleMap {
		if t.Prefix != "" {
			key = path.Join(path.Clean(t.Prefix), key)
		}
		if !t.SecretAsKey {
			// Encode as json
			payload, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("unable to encode value as JSON for %q: %w", key, err)
			}

			// Insert in KV store.
			if err := t.Store.Put(ctx, key, payload); err != nil {
				return fmt.Errorf("unable to publish %q secret in store: %w", key, err)
			}
		} else {
			// Range over secrets
			secrets, ok := value.(bundle.KV)
			if !ok {
				continue
			}

			// Publish each secret as a leaf.
			for secKey, secValue := range secrets {
				// Insert in KV store.
				if err := t.Store.Put(ctx, path.Join(key, secKey), []byte(fmt.Sprintf("%v", secValue))); err != nil {
					return fmt.Errorf("unable to publish %q secret in store: %w", key, err)
				}
			}
		}
	}

	// No error
	return nil
}
