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
	"zntr.io/harp/v2/pkg/tasks/from"
)

// -----------------------------------------------------------------------------

var fromDumpCmd = func() *cobra.Command {
	var (
		inputPath  string
		outputPath string
	)
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Import from bundle dump output as a secret container",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-from-dump", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &from.BundleDumpTask{
				JSONReader:   cmdutil.FileReader(inputPath),
				OutputWriter: cmdutil.FileWriter(outputPath),
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "Bundle dump input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "-", "Container output ('-' for stdout or filename)")

	return cmd
}
