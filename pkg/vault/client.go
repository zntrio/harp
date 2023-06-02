// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"

	"zntr.io/harp/v2/pkg/vault/cubbyhole"
	"zntr.io/harp/v2/pkg/vault/kv"
	"zntr.io/harp/v2/pkg/vault/transit"
)

// -----------------------------------------------------------------------------

// ServiceFactory defines Vault client cervice contract.
type ServiceFactory interface {
	KV(mountPath string) (kv.Service, error)
	Transit(mounthPath, keyName string) (transit.Service, error)
	Cubbyhole(mountPath string) (cubbyhole.Service, error)
}

// -----------------------------------------------------------------------------

// DefaultClient initialize a Vault client and wrap it in a Service factory.
func DefaultClient() (ServiceFactory, error) {
	// Initialize default config
	conf := api.DefaultConfig()

	// Initialize vault client
	vaultClient, err := api.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	// Delegate to other constructor.
	return FromVaultClient(vaultClient)
}

// FromVaultClient wraps an existing Vault client as a Service factory.
func FromVaultClient(vaultClient *api.Client) (ServiceFactory, error) {
	// Return wrapped client.
	return &client{
		Client: vaultClient,
	}, nil
}

// -----------------------------------------------------------------------------

// Client wrpas original Vault client instance to provide service factory.
type client struct {
	*api.Client
}

func (c *client) KV(mountPath string) (kv.Service, error) {
	return kv.New(c.Client, mountPath)
}

func (c *client) Transit(mountPath, keyName string) (transit.Service, error) {
	return transit.New(c.Client, mountPath, keyName)
}

func (c *client) Cubbyhole(mountPath string) (cubbyhole.Service, error) {
	return cubbyhole.New(c.Client, mountPath)
}
