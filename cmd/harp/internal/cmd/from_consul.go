// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"context"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/kv/consul"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/from"
)

// -----------------------------------------------------------------------------

type fromConsulParams struct {
	outputPath           string
	basePaths            []string
	lastPathItemAsSecret bool
}

var fromConsulCmd = func() *cobra.Command {
	var params fromConsulParams

	cmd := &cobra.Command{
		Use:   "consul",
		Short: "Extract KV pairs from Hashicorp Consul KV Store",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-kv-from-consul", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			runFromConsul(ctx, &params)
		},
	}

	// Add parameters
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Container output path ('-' for stdout)")
	cmd.Flags().StringSliceVar(&params.basePaths, "paths", []string{}, "Exported base paths")
	cmd.Flags().BoolVarP(&params.lastPathItemAsSecret, "last-path-item-as-secret-key", "k", false, "Use the last path element as secret key")

	return cmd
}

func runFromConsul(ctx context.Context, params *fromConsulParams) {
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
	t := &from.ExtractKVTask{
		Store:                   store,
		ContainerWriter:         cmdutil.FileWriter(params.outputPath),
		BasePaths:               params.basePaths,
		LastPathItemAsSecretKey: params.lastPathItemAsSecret,
	}

	// Run the task
	if err := t.Run(ctx); err != nil {
		log.For(ctx).Fatal("unable to execute kv extract task", zap.Error(err))
		return
	}
}
