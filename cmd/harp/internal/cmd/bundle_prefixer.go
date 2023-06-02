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
	"zntr.io/harp/v2/pkg/tasks/bundle"
)

// -----------------------------------------------------------------------------.
type bundlePrefixerParams struct {
	inputPath  string
	outputPath string
	prefix     string
	remove     bool
}

var bundlePrefixerCmd = func() *cobra.Command {
	params := bundlePrefixerParams{}

	cmd := &cobra.Command{
		Use:   "prefixer",
		Short: "Simple package prefix operaton",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-prefixer", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.PrefixerTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				OutputWriter:    cmdutil.FileWriter(params.outputPath),
				Prefix:          params.prefix,
				Remove:          params.remove,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "", "Container output ('-' for stdout or a filename)")
	cmd.Flags().StringVar(&params.prefix, "prefix", "", "Specify prefix to prepend")
	cmd.Flags().BoolVarP(&params.remove, "remove", "r", false, "Remove the given prefix from the package paths")

	return cmd
}
