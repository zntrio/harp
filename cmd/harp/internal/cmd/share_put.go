// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/share"
)

// -----------------------------------------------------------------------------

var sharePutCmd = func() *cobra.Command {
	var (
		inputPath     string
		backendPrefix string
		namespace     string
		ttl           time.Duration
		jsonOutput    bool
	)

	cmd := &cobra.Command{
		Use:   "put",
		Short: "Put secret in Vault Cubbyhole and return a wrapped token",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "share-put", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &share.PutTask{
				InputReader:    cmdutil.FileReader(inputPath),
				OutputWriter:   cmdutil.StdoutWriter(),
				BackendPrefix:  backendPrefix,
				VaultNamespace: namespace,
				TTL:            ttl,
				JSONOutput:     jsonOutput,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&backendPrefix, "prefix", "cubbyhole", "Vault backend prefix")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Vault namespace")
	cmd.Flags().DurationVar(&ttl, "ttl", 30*time.Second, "Token expiration")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Display result as json")

	return cmd
}
