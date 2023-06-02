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
	"zntr.io/harp/v2/pkg/tasks/template"
)

type valuesParams struct {
	OutputPath   string
	ValueFiles   []string
	Values       []string
	StringValues []string
	FileValues   []string
}

// -----------------------------------------------------------------------------

var valuesCmd = func() *cobra.Command {
	params := &valuesParams{}

	cmd := &cobra.Command{
		Use:     "values",
		Aliases: []string{"v"},
		Short:   "Template value preprocessor",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-values", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &template.ValueTask{
				OutputWriter: cmdutil.FileWriter(params.OutputPath),
				ValueFiles:   params.ValueFiles,
				Values:       params.Values,
				StringValues: params.StringValues,
				FileValues:   params.FileValues,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.OutputPath, "out", "", "Output file ('-' for stdout or a filename)")
	cmd.Flags().StringArrayVarP(&params.ValueFiles, "values", "f", []string{}, "Specifies value files to load. Use <path>:<type>[:<prefix>] to override type detection (json,yaml,xml,hocon,toml)")
	cmd.Flags().StringArrayVar(&params.Values, "set", []string{}, "Specifies value (k=v)")
	cmd.Flags().StringArrayVar(&params.StringValues, "set-string", []string{}, "Specifies value (k=string)")
	cmd.Flags().StringArrayVar(&params.FileValues, "set-file", []string{}, "Specifies value (k=filepath)")

	return cmd
}
