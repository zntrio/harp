// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/kv/consul"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/to"
)

// -----------------------------------------------------------------------------

type toConsulParams struct {
	inputPath    string
	secretAsLeaf bool
	prefix       string
}

var toConsulCmd = func() *cobra.Command {
	var params toConsulParams

	cmd := &cobra.Command{
		Use:   "consul",
		Short: "Publish bundle data into HashiCorp Consul",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-kv-to-consul", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Create Consul client config from environment.
			config := api.DefaultConfig()

			// Creates a new client
			client, err := api.NewClient(config)
			if err != nil {
				log.For(ctx).Fatal("unable to connect to consul cluster", zap.Error(err))
				return
			}

			// Prepare store.
			store := consul.Store(client.KV())
			defer log.SafeClose(store, "unable to close consul store")

			// Delegate to task
			t := &to.PublishKVTask{
				Store:           store,
				ContainerReader: cmdutil.FileReader(params.inputPath),
				SecretAsKey:     params.secretAsLeaf,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute kv extract task", zap.Error(err))
				return
			}
		},
	}

	// Add parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Container path ('-' for stdin or filename)")
	cmd.Flags().BoolVarP(&params.secretAsLeaf, "secret-as-leaf", "s", false, "Expand package path to secrets for provisioning")
	cmd.Flags().StringVar(&params.prefix, "prefix", "", "Path prefix for insertion")

	return cmd
}
