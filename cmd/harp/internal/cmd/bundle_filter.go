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
type bundleFilterParams struct {
	inputPath      string
	outputPath     string
	excludePaths   []string
	keepPaths      []string
	jmesPath       string
	regoPolicy     string
	celExpressions []string
	reverseLogic   bool
}

var bundleFilterCmd = func() *cobra.Command {
	params := &bundleFilterParams{}

	longDesc := cmdutil.LongDesc(`
	Create a new Bundle based on applied package matchers.

	Filtering a Bundle consists in reducing the Bundle packages using a matcher
	applied on the Bundle and Package model to select them, and export them
	in another Bundle.

	In order to filter packages, you can use :
	* a package name selector
	* a JMES query
	* a REGO policy
	* a Set of CEL expressions

	Bundle package filtering capabilities are the root of the secret management
	by contract. Filter commands can be pipelined to produce complex filtering
	pipelines and target the appropriate secrets.

	TIP: Use this command to debug your BundlePatch matchers.`)

	examples := cmdutil.Examples(`
	# Exclude specific packages by name from STDIN bundle to STDOUT.
	harp bundle filter --exclude "$app/(staging|production)/*"

	# Exclude specific packages by name from file bundle to STDOUT
	harp bundle filter --in customer.bundle --exclude "$app/(staging|production)/*"

	# Keep specific packages by name
	harp bundle filter --keep "$app/(staging|production)/*"

	# Filter packages using a JMES query (context is the package)
	harp bundle filter --query "labels.deprecated == 'true'"

	# Filter packages using a JMES query (context is the package) to a file based bundle.
	harp bundle filter --query "labels.deprecated == 'true'" --out deprecated.bundle

	# Filter packages using a REGO policy
	harp bundle filter --policy deprecated.rego

	# Filter packages using a CEL matcher expressions (associated with AND logic if multiple)
	harp bundle filter --cel "p.match_secret('*Key')"

	# Reverse the matcher logic
	harp bundle filter --not <matcher>`)

	cmd := &cobra.Command{
		Use:     "filter",
		Aliases: []string{"grep", "f"},
		Short:   "Filter package names",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-filter", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.FilterTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				OutputWriter:    cmdutil.FileWriter(params.outputPath),
				ExcludePaths:    params.excludePaths,
				KeepPaths:       params.keepPaths,
				JMESPath:        params.jmesPath,
				RegoPolicy:      params.regoPolicy,
				CELExpressions:  params.celExpressions,
				ReverseLogic:    params.reverseLogic,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "", "Container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "", "Container path ('-' for stdout or filename)")
	cmd.Flags().StringArrayVar(&params.excludePaths, "exclude", []string{}, "Exclude path")
	cmd.Flags().StringArrayVar(&params.keepPaths, "keep", []string{}, "Keep path")
	cmd.Flags().StringVar(&params.jmesPath, "query", "", "JMESPath query used as package filter")
	cmd.Flags().StringVar(&params.regoPolicy, "policy", "", "OPA Rego policy file as package filter")
	cmd.Flags().StringArrayVar(&params.celExpressions, "cel", []string{}, "CEL expression as package filter (multiple)")
	cmd.Flags().BoolVar(&params.reverseLogic, "not", false, "Reverse filter logic expression")

	return cmd
}
