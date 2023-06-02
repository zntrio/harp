// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// SecretGetter pull a secret from Vault using given path.
//
// To be used of template function.
func SecretGetter(ctx context.Context, client *api.Client) func(string) (map[string]interface{}, error) {
	return func(path string) (map[string]interface{}, error) {
		// Create dedicated service reader
		service, err := New(client, path, WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("unable to prepare vault reader for path %q: %w", path, err)
		}

		// Delegate to reader
		data, _, err := service.Read(ctx, path)
		return data, err
	}
}
