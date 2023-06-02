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

var toObjectCmd = func() *cobra.Command {
	var (
		inputPath  string
		outputPath string
		expand     bool
		outputType string
	)

	cmd := &cobra.Command{
		Use:   "object",
		Short: "Export all data of a secret container as JSON / YAML / TOML.",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-to-object", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &to.ObjectTask{
				ContainerReader: cmdutil.FileReader(inputPath),
				OutputWriter:    cmdutil.FileWriter(outputPath),
				Expand:          expand,
				TOML:            outputType == "toml",
				YAML:            outputType == "yaml",
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "Container path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "-", "Container output ('-' for stdout or filename)")
	cmd.Flags().StringVar(&outputType, "format", "json", "Output format (json / yaml / toml)")
	cmd.Flags().BoolVar(&expand, "expand", false, "Expand package paths as embedded maps")

	return cmd
}
