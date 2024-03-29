// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Consul creates a test consul server inside a Docker container.
// nolint: contextcheck // false positive
func Consul(_ context.Context, tb testing.TB) string {
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("couldn't connect to docker: %v", err)
		return ""
	}
	pool.MaxWait = 10 * time.Second

	// Prepare bootstrap configuration
	config := struct {
		Datacenter       string `json:"datacenter,omitempty"`
		ACLDatacenter    string `json:"acl_datacenter,omitempty"`
		ACLDefaultPolicy string `json:"acl_default_policy,omitempty"`
		ACLMasterToken   string `json:"acl_master_token,omitempty"`
	}{
		Datacenter:       "test",
		ACLDatacenter:    "test",
		ACLDefaultPolicy: "deny",
		ACLMasterToken:   "test",
	}

	// Encode configuration as JSON
	encodedConfig, errConfig := json.Marshal(config)
	if errConfig != nil {
		tb.Fatalf("couldn't serialize configuration as json: %v", errConfig)
		return ""
	}

	// Start zookeeper server
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "consul",
		Tag:        "1.10.3",
		Cmd:        []string{"agent", "-dev", "-client", "0.0.0.0"},
		Env:        []string{fmt.Sprintf("CONSUL_LOCAL_CONFIG=%s", encodedConfig)},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		tb.Fatalf("couldn't start resource: %v", err)
		return ""
	}

	// Set expiration
	if err := resource.Expire(15 * 60); err != nil {
		tb.Error("unable to set expiration value for the container")
	}

	// Cleanup function
	tb.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			tb.Errorf("couldn't purge container: %v", err)
			return
		}
	})

	consulURI := fmt.Sprintf("localhost:%s", resource.GetPort("8500/tcp"))

	// Wait until connection is ready
	if err := pool.Retry(func() (err error) {
		config := api.DefaultConfig()
		config.Address = consulURI
		config.Token = "test"

		// Create client instance.
		client, err := api.NewClient(config)
		if err != nil {
			return fmt.Errorf("unable to connect to the server: %w", err)
		}

		// Try to write data.
		_, err = client.KV().Put(&api.KVPair{
			Key:   "ready",
			Value: []byte("ready"),
		}, nil)

		// Check connection state
		return err
	}); err != nil {
		tb.Fatalf("zk server never ready: %v", err)
		return ""
	}

	// Return connection uri
	return consulURI
}
