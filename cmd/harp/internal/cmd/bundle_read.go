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

// -----------------------------------------------------------------------------

var bundleReadCmd = func() *cobra.Command {
	var (
		inputPath   string
		packageName string
		secretKey   string
	)

	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read a secret from bundle",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-read", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.ReadTask{
				ContainerReader: cmdutil.FileReader(inputPath),
				OutputWriter:    cmdutil.StdoutWriter(),
				PackageName:     packageName,
				SecretKey:       secretKey,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "", "Container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&packageName, "path", "", "Secret path")
	log.CheckErr("unable to mark 'path' flag as required.", cmd.MarkFlagRequired("path"))
	cmd.Flags().StringVar(&secretKey, "field", "", "Secret field")

	return cmd
}
