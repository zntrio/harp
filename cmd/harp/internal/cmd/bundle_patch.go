// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/bundle/patch"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/bundle"
	tplcmdutil "zntr.io/harp/v2/pkg/template/cmdutil"
)

// -----------------------------------------------------------------------------.
type bundlePatchParams struct {
	inputPath         string
	outputPath        string
	patchPath         string
	valueFiles        []string
	values            []string
	stringValues      []string
	fileValues        []string
	stopAtRuleIndex   int
	stopAtRuleID      string
	ignoreRuleIDs     []string
	ignoreRuleIndexes []int
}

var bundlePatchCmd = func() *cobra.Command {
	params := &bundlePatchParams{}

	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Apply patch to the given bundle",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-patch", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Load values
			valueOpts := tplcmdutil.ValueOptions{
				ValueFiles:   params.valueFiles,
				Values:       params.values,
				StringValues: params.stringValues,
				FileValues:   params.fileValues,
			}
			values, err := valueOpts.MergeValues()
			if err != nil {
				log.For(ctx).Fatal("unable to process values", zap.Error(err))
			}

			// Prepare patch options.
			opts := []patch.OptionFunc{
				patch.WithStopAtRuleID(params.stopAtRuleID),
				patch.WithStopAtRuleIndex(params.stopAtRuleIndex),
				patch.WithIgnoreRuleIDs(params.ignoreRuleIDs...),
				patch.WithIgnoreRuleIndexes(params.ignoreRuleIndexes...),
			}

			// Prepare task
			t := &bundle.PatchTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				PatchReader:     cmdutil.FileReader(params.patchPath),
				OutputWriter:    cmdutil.FileWriter(params.outputPath),
				Values:          values,
				Options:         opts,
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
	cmd.Flags().StringVar(&params.patchPath, "spec", "", "Patch specification path ('-' for stdin or filename)")
	log.CheckErr("unable to mark 'spec' flag as required.", cmd.MarkFlagRequired("spec"))
	cmd.Flags().StringArrayVar(&params.valueFiles, "values", []string{}, "Specifies value files to load")
	cmd.Flags().StringArrayVar(&params.values, "set", []string{}, "Specifies value (k=v)")
	cmd.Flags().StringArrayVar(&params.stringValues, "set-string", []string{}, "Specifies value (k=string)")
	cmd.Flags().StringArrayVar(&params.fileValues, "set-file", []string{}, "Specifies value (k=filepath)")
	cmd.Flags().StringVar(&params.stopAtRuleID, "stop-at-rule-id", "", "Stop patch evaluation before the given rule ID")
	cmd.Flags().IntVar(&params.stopAtRuleIndex, "stop-at-rule-index", -1, "Stop patch evaluation before the given rule index (0 for first rule)")
	cmd.Flags().StringArrayVar(&params.ignoreRuleIDs, "ignore-rule-id", []string{}, "List of Rule identifier to ignore during evaluation")
	cmd.Flags().IntSliceVar(&params.ignoreRuleIndexes, "ignore-rule-index", []int{}, "List of Rule index to ignore during evaluation")

	return cmd
}
