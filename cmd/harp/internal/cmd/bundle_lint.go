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
type bundleLintParams struct {
	inputPath string
	specPath  string
}

var bundleLintCmd = func() *cobra.Command {
	params := &bundleLintParams{}

	longDesc := cmdutil.LongDesc(`
	Apply a RuleSet specification to the given bundle.

	This command is used to check a Bundle structure (Package => Secrets).
	A control gate could be implemented with this command to enforce a bundle
	structure by decoupling the bundle content and the usage contract.`)

	examples := cmdutil.Examples(`
	# Lint a bundle from STDIN
	harp bundle lint --spec cso.yaml`)

	cmd := &cobra.Command{
		Use:     "lint",
		Short:   "Lint the bundle using the given RuleSet spec",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-lint", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.LintTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				RuleSetReader:   cmdutil.FileReader(params.specPath),
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.specPath, "spec", "", "RuleSet specification path ('-' for stdin or filename)")
	log.CheckErr("unable to mark 'spec' flag as required.", cmd.MarkFlagRequired("spec"))

	return cmd
}
