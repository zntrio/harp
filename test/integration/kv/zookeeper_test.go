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

	"github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"

	"zntr.io/harp/v2/pkg/kv/zookeeper"
	"zntr.io/harp/v2/test/integration/resource"
)

// -----------------------------------------------------------------------------

func TestWithZookeeper(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create zk instance
	kvURI := resource.Zookeeper(ctx, t)

	// Create zk client
	conn, _, err := zk.Connect([]string{kvURI}, 10*time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// Initialize KV Store
	s := zookeeper.Store(conn)

	// Run test suite
	t.Run("store", testSuite(ctx, s))
}
