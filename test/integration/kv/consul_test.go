// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build integration

package kv

import (
	"context"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/kv/consul"
	"zntr.io/harp/v2/test/integration/resource"
)

// -----------------------------------------------------------------------------

func TestWithConsul(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create zk instance
	kvURI := resource.Consul(ctx, t)

	config := api.DefaultConfig()
	config.Address = kvURI
	config.Token = "test"

	// Create client instance.
	client, err := api.NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Initialize KV Store
	s := consul.Store(client.KV())

	// Run test suite
	t.Run("store", testSuite(ctx, s))
}
