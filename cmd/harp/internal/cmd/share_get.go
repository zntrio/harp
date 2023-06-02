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
	"zntr.io/harp/v2/pkg/tasks/share"
)

// -----------------------------------------------------------------------------

var shareGetCmd = func() *cobra.Command {
	var (
		outputPath    string
		backendPrefix string
		namespace     string
		token         string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get wrapped secret from Vault",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "share-get", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &share.GetTask{
				OutputWriter:   cmdutil.FileWriter(outputPath),
				BackendPrefix:  backendPrefix,
				VaultNamespace: namespace,
				Token:          token,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&backendPrefix, "prefix", "cubbyhole", "Vault backend prefix")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Vault namespace")
	cmd.Flags().StringVar(&token, "token", "", "Wrapped token")
	log.CheckErr("unable to mark 'token' flag as required.", cmd.MarkFlagRequired("token"))

	return cmd
}
