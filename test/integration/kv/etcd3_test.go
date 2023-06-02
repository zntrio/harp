// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build integration

package kv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"

	"zntr.io/harp/v2/pkg/kv/etcd3"
	"zntr.io/harp/v2/test/integration/resource"
)

// -----------------------------------------------------------------------------

func TestWithEtcd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create zk instance
	kvURI := resource.Etcd(ctx, t)

	// Create zk client
	client, errClient := clientv3.New(clientv3.Config{
		Endpoints:   []string{kvURI},
		DialTimeout: 5 * time.Second,
	})
	assert.NoError(t, errClient)
	assert.NotNil(t, client)

	// Initialize KV Store
	s := etcd3.Store(client)

	// Run test suite
	t.Run("store", testSuite(ctx, s))
}
