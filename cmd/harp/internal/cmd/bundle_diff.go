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
type bundleDiffParams struct {
	sourcePath      string
	destinationPath string
	generatePatch   bool
	outputPath      string
}

var bundleDiffCmd = func() *cobra.Command {
	params := &bundleDiffParams{}

	longDesc := cmdutil.LongDesc(`
	Compute Bundle object differences.

	Useful to debug a BundlePatch application and watch for a Bundle alteration.
	`)

	examples := cmdutil.Examples(`
	# Diff a bundle from STD and a file based one
	harp bundle diff --old - --new rotated.bundle

	# Generate a BundlePatch from differences
	harp bundle diff --old - --new rotated.bundle --patch --out rotation.yaml`)

	cmd := &cobra.Command{
		Use:     "diff",
		Short:   "Display bundle differences",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-diff", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.DiffTask{
				SourceReader:      cmdutil.FileReader(params.sourcePath),
				DestinationReader: cmdutil.FileReader(params.destinationPath),
				OutputWriter:      cmdutil.FileWriter(params.outputPath),
				GeneratePatch:     params.generatePatch,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.sourcePath, "old", "", "Container path ('-' for stdin or filename)")
	log.CheckErr("unable to mark 'old' flag as required.", cmd.MarkFlagRequired("old"))
	cmd.Flags().StringVar(&params.destinationPath, "new", "", "Container path ('-' for stdin or filename)")
	log.CheckErr("unable to mark 'new' flag as required.", cmd.MarkFlagRequired("new"))
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output ('-' for stdout or filename)")
	cmd.Flags().BoolVar(&params.generatePatch, "patch", false, "Output as a bundle patch")

	return cmd
}
