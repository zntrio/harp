// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package consul

import (
	"context"
	"errors"
	"fmt"
	"strings"

	api "github.com/hashicorp/consul/api"
	"zntr.io/harp/v2/pkg/kv"
	"zntr.io/harp/v2/pkg/sdk/types"
)

type consulDriver struct {
	client Client
}

func Store(client Client) kv.Store {
	return &consulDriver{
		client: client,
	}
}

// -----------------------------------------------------------------------------

func (d *consulDriver) Get(_ context.Context, key string) (*kv.Pair, error) {
	// Check arguments
	if types.IsNil(d.client) {
		return nil, errors.New("consul: unable to query with nil client")
	}

	// Retrieve from backend
	item, meta, err := d.client.Get(d.normalize(key), &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	})
	if err != nil {
		return nil, fmt.Errorf("consul: unable to retrieve %q key: %w", key, err)
	}
	if item == nil {
		return nil, kv.ErrKeyNotFound
	}

	// No error
	return &kv.Pair{
		Key:     item.Key,
		Value:   item.Value,
		Version: meta.LastIndex,
	}, nil
}

func (d *consulDriver) Put(_ context.Context, key string, value []byte) error {
	// Check arguments
	if types.IsNil(d.client) {
		return errors.New("consul: unable to query with nil client")
	}

	// Prepare the item to put
	item := &api.KVPair{
		Key:   d.normalize(key),
		Value: value,
	}

	// Delegate to client
	if _, err := d.client.Put(item, nil); err != nil {
		return fmt.Errorf("consul: unable to put %q value: %w", key, err)
	}

	// No error
	return nil
}

func (d *consulDriver) Delete(ctx context.Context, key string) error {
	// Check arguments
	if types.IsNil(d.client) {
		return errors.New("consul: unable to query with nil client")
	}

	// Retrieve from store
	found, err := d.Exists(ctx, key)
	if err != nil {
		return fmt.Errorf("consul: unable to retrieve %q for deletion: %w", key, err)
	}
	if !found {
		return kv.ErrKeyNotFound
	}

	// Delete the value
	if _, err := d.client.Delete(d.normalize(key), nil); err != nil {
		return fmt.Errorf("consul: unable to delete %q: %w", key, err)
	}

	// No error
	return nil
}

func (d *consulDriver) Exists(ctx context.Context, key string) (bool, error) {
	// Retrieve from stroe
	_, err := d.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrKeyNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("consul: unable to check key %q existence: %w", key, err)
	}

	// No error
	return true, nil
}

func (d *consulDriver) List(_ context.Context, basePath string) ([]*kv.Pair, error) {
	// Check arguments
	if types.IsNil(d.client) {
		return nil, errors.New("consul: unable to query with nil client")
	}

	// List keys from base path
	items, _, err := d.client.List(d.normalize(basePath), nil)
	if err != nil {
		return nil, fmt.Errorf("consul: unable to list keys from %q: %w", basePath, err)
	}
	if len(items) == 0 {
		return nil, kv.ErrKeyNotFound
	}

	// Unpack values
	results := []*kv.Pair{}
	for _, item := range items {
		// Skip first item as base path
		if item.Key == basePath {
			continue
		}
		results = append(results, &kv.Pair{
			Key:     item.Key,
			Value:   item.Value,
			Version: item.ModifyIndex,
		})
	}

	// No error
	return results, nil
}

func (d *consulDriver) Close() error {
	// No error
	return nil
}

// -----------------------------------------------------------------------------

// Normalize the key for usage in Consul.
func (d *consulDriver) normalize(key string) string {
	key = kv.Normalize(key)
	return strings.TrimPrefix(key, "/")
}
