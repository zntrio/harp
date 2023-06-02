// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package from

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle"
	"zntr.io/harp/v2/pkg/bundle/secret"
	"zntr.io/harp/v2/pkg/kv"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks"
)

type ExtractKVTask struct {
	_                       struct{}
	ContainerWriter         tasks.WriterProvider
	BasePaths               []string
	Store                   kv.Store
	LastPathItemAsSecretKey bool
}

func (t *ExtractKVTask) Run(ctx context.Context) error {
	packages := map[string]*bundlev1.Package{}

	// For each base path
	for _, basePath := range t.BasePaths {
		// List recusively items
		items, err := t.Store.List(ctx, basePath)
		if err != nil {
			return fmt.Errorf("unable to extract key from store: %w", err)
		}

		// Prepare a package using each item
		for _, item := range items {
			// Extract packageName
			packageName := t.extractPackageName(item.Key)

			// Check if packages is already instancied
			p, ok := packages[packageName]
			if !ok {
				p = &bundlev1.Package{
					Name: packageName,
					Secrets: &bundlev1.SecretChain{
						Version:         uint32(0),
						Data:            make([]*bundlev1.KV, 0),
						NextVersion:     nil,
						PreviousVersion: nil,
					},
				}
			}

			// Try to extract value as a json map
			var secretData map[string]interface{}
			errJSON := json.Unmarshal(item.Value, &secretData)
			if errJSON != nil {
				log.For(ctx).Debug("data could not be decoded as json", zap.Error(errJSON))

				// Create an arbitrary secret key
				secretKey := strings.TrimPrefix(strings.TrimPrefix(item.Key, kv.GetDirectory(item.Key)), "/")

				log.For(ctx).Debug("Creating secret for package", zap.String("package", packageName), zap.String("secret", secretKey))

				// Pack secret value
				s, errPack := t.packSecret(secretKey, string(item.Value))
				if errPack != nil {
					return fmt.Errorf("unable to pack secret value for path %q with key %q : %w", item.Key, secretKey, errPack)
				}

				// Add secret to package
				p.Secrets.Data = append(p.Secrets.Data, s)
			} else {
				// Iterate over secret bundle
				for k, v := range secretData {
					// Pack secret value
					s, errPack := t.packSecret(k, v)
					if errPack != nil {
						return fmt.Errorf("unable to pack secret value for path %q with key %q : %w", item.Key, k, errPack)
					}

					// Add secret to package
					p.Secrets.Data = append(p.Secrets.Data, s)
				}
			}

			// Update package map
			packages[packageName] = p
		}
	}

	// Prepare a bundle
	b := &bundlev1.Bundle{
		Packages: make([]*bundlev1.Package, 0),
	}

	// Copy package to bundle
	for _, p := range packages {
		b.Packages = append(b.Packages, p)
	}

	// Create container
	writer, err := t.ContainerWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize container writer: %w", err)
	}

	// Dump bundle
	if err = bundle.ToContainerWriter(writer, b); err != nil {
		return fmt.Errorf("unable to produce exported bundle: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------.
func (t *ExtractKVTask) packSecret(key string, value interface{}) (*bundlev1.KV, error) {
	// Pack secret value
	payload, err := secret.Pack(value)
	if err != nil {
		return nil, fmt.Errorf("unable to pack secret %q: %w", key, err)
	}

	// Build the secret object
	return &bundlev1.KV{
		Key:   key,
		Type:  fmt.Sprintf("%T", value),
		Value: payload,
	}, nil
}

func (t *ExtractKVTask) extractPackageName(key string) string {
	if !t.LastPathItemAsSecretKey {
		return strings.TrimPrefix(strings.TrimSuffix(key, "/"), "/")
	}

	// Extract directory
	return strings.TrimPrefix(kv.GetDirectory(key), "/")
}
