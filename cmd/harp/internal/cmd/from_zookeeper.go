// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"context"
	"time"

	zk "github.com/go-zookeeper/zk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/kv/zookeeper"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/from"
)

// -----------------------------------------------------------------------------

type fromZookeeperParams struct {
	outputPath           string
	basePaths            []string
	lastPathItemAsSecret bool

	endpoints   []string
	dialTimeout time.Duration
}

var fromZookeeperCmd = func() *cobra.Command {
	var params fromZookeeperParams
	cmd := &cobra.Command{
		Use:     "zookeeper",
		Aliases: []string{"zk"},
		Short:   "Extract KV pairs from Apache Zookeeper KV Store",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-kv-from-zookeeper", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			runFromZookeeper(ctx, &params)
		},
	}

	// Add parameters
	cmd.Flags().StringArrayVar(&params.endpoints, "endpoints", []string{"127.0.0.1:2181"}, "Zookeeper client endpoints")
	cmd.Flags().DurationVar(&params.dialTimeout, "dial-timeout", 15*time.Second, "Zookeeper client dial timeout")
	cmd.Flags().BoolVarP(&params.lastPathItemAsSecret, "last-path-item-as-secret-key", "k", false, "Use the last path element as secret key")

	return cmd
}

func runFromZookeeper(ctx context.Context, params *fromZookeeperParams) {
	// Create config
	// nolint: contextcheck // zk lib doesn't support to pass a caller context yet
	client, _, err := zk.Connect(params.endpoints, params.dialTimeout)
	if err != nil {
		log.For(ctx).Fatal("unable to connect to zookeeper cluster", zap.Error(err))
		return
	}

	// Prepare store.
	store := zookeeper.Store(client)
	defer log.SafeClose(store, "unable to close zk store")

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
