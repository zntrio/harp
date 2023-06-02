// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"time"

	zk "github.com/go-zookeeper/zk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/kv/zookeeper"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/to"
)

// -----------------------------------------------------------------------------

type toZookeeperParams struct {
	inputPath    string
	secretAsLeaf bool
	prefix       string

	endpoints   []string
	dialTimeout time.Duration
}

var toZookeeperCmd = func() *cobra.Command {
	var params toZookeeperParams

	cmd := &cobra.Command{
		Use:     "zookeeper",
		Aliases: []string{"zk"},
		Short:   "Publish bundle data into Apache Zookeeper",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-kv-to-zookeeper", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Create config
			//nolint: contextcheck // zk lib doesn't support to pass a caller context yet
			client, _, err := zk.Connect(params.endpoints, params.dialTimeout)
			if err != nil {
				log.For(ctx).Fatal("unable to connect to zookeeper cluster", zap.Error(err))
				return
			}

			// Prepare store.
			store := zookeeper.Store(client)
			defer log.SafeClose(store, "unable to close zk store")

			// Delegate to task
			t := &to.PublishKVTask{
				Store:           store,
				ContainerReader: cmdutil.FileReader(params.inputPath),
				SecretAsKey:     params.secretAsLeaf,
				Prefix:          params.prefix,
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

	cmd.Flags().StringArrayVar(&params.endpoints, "endpoints", []string{"127.0.0.1:2181"}, "Zookeeper client endpoints")
	cmd.Flags().DurationVar(&params.dialTimeout, "dial-timeout", 15*time.Second, "Zookeeper client dial timeout")

	return cmd
}
