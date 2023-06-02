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

var fromObjectCmd = func() *cobra.Command {
	var (
		inputPath  string
		outputPath string
		inputType  string
	)
	cmd := &cobra.Command{
		Use:   "object",
		Short: "Convert an object to a secret-container",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-from-object", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &from.ObjectTask{
				ObjectReader: cmdutil.FileReader(inputPath),
				OutputWriter: cmdutil.FileWriter(outputPath),
				JSON:         inputType == "json",
				YAML:         inputType == "yaml",
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "YAML object ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "-", "Container output ('-' for stdout or filename)")
	cmd.Flags().StringVar(&inputType, "format", "yaml", "Input file format (yaml or json)")

	return cmd
}
