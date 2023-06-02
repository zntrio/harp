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
type bundleDumpParams struct {
	inputPath      string
	dataOnly       bool
	metadataOnly   bool
	pathOnly       bool
	jmesPathFilter string
	skipTemplate   bool
}

var bundleDumpCmd = func() *cobra.Command {
	params := &bundleDumpParams{}

	longDesc := cmdutil.LongDesc(`
	Inspect a Bundle object.

	Harp Bundles is a structure designed to hold additional properties associated
	to a path (package name) and values (secrets). For your pipeline usages, you
	can store annotations, labels and user data which can be consumed and/or
	produced during the secret management pipeline execution.

	The Bundle object specification can be consulted here -	https://ela.st/harp-spec-bundle
	`)

	examples := cmdutil.Examples(`
	# Dump a JSON representation of a Bundle object from STDIN
	harp bundle dump

	# Dump a JSON map containing package name as key and associated secret kv
	harp bundle dump --data-only

	# Dump a JSON map containing package name as key and associated metadata
	harp bundle dump --metadata-only

	# Dump all package paths as a list (useful for xargs usage)
	harp bundle dump --path-only

	# Dump a Bundle using a JMEFilter query
	harp bundle dump --query <jmesfilter query>

	# Dump a bundle content excluding the template used to generate
	harp bundle dump --skip-template`)

	cmd := &cobra.Command{
		Use:     "dump",
		Short:   "Dump as JSON",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-dump", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.DumpTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				OutputWriter:    cmdutil.StdoutWriter(),
				DataOnly:        params.dataOnly,
				MetadataOnly:    params.metadataOnly,
				PathOnly:        params.pathOnly,
				JMESPathFilter:  params.jmesPathFilter,
				IgnoreTemplate:  params.skipTemplate,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "", "Container input ('-' for stdin or filename)")
	cmd.Flags().BoolVar(&params.dataOnly, "content-only", false, "Display content only (data-only alias)")
	cmd.Flags().BoolVar(&params.dataOnly, "data-only", false, "Display data only")
	cmd.Flags().BoolVar(&params.metadataOnly, "metadata-only", false, "Display metadata only")
	cmd.Flags().BoolVar(&params.pathOnly, "path-only", false, "Display path only")
	cmd.Flags().StringVar(&params.jmesPathFilter, "query", "", "Specify a JMESPath query to format output")
	cmd.Flags().BoolVar(&params.skipTemplate, "skip-template", false, "Drop template from dump")

	return cmd
}
