// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/to"
)

// -----------------------------------------------------------------------------

var toVaultCmd = func() *cobra.Command {
	var (
		inputPath         string
		backendPrefix     string
		namespace         string
		withMetadata      bool
		withVaultMetadata bool
		maxWorkerCount    int64
	)

	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Push a secret container in Hashicorp Vault",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-to-vault", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &to.VaultTask{
				ContainerReader: cmdutil.FileReader(inputPath),
				BackendPrefix:   backendPrefix,
				PushMetadata:    withMetadata || withVaultMetadata,
				AsVaultMetadata: withVaultMetadata,
				VaultNamespace:  namespace,
				MaxWorkerCount:  maxWorkerCount,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "Container path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&backendPrefix, "prefix", "", "Vault backend prefix")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Vault namespace")
	cmd.Flags().BoolVar(&withMetadata, "with-metadata", false, "Push container metadata as secret data")
	cmd.Flags().BoolVar(&withVaultMetadata, "with-vault-metadata", false, "Push container metadata as secret metadata (requires Vault >=1.9)")
	cmd.Flags().Int64Var(&maxWorkerCount, "worker-count", 4, "Active worker count limit")

	return cmd
}
