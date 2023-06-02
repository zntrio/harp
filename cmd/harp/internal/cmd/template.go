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

type templateParams struct {
	InputPath     string
	OutputPath    string
	ValueFiles    []string
	SecretLoaders []string
	Values        []string
	StringValues  []string
	FileValues    []string
	LeftDelims    string
	RightDelims   string
	AltDelims     bool
	RootPath      string
}

// -----------------------------------------------------------------------------

var templateCmd = func() *cobra.Command {
	params := &templateParams{}

	cmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"t", "tpl"},
		Short:   "Read a template and execute it",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "template-render", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &template.RenderTask{
				InputReader:   cmdutil.FileReader(params.InputPath),
				OutputWriter:  cmdutil.FileWriter(params.OutputPath),
				ValueFiles:    params.ValueFiles,
				Values:        params.Values,
				StringValues:  params.StringValues,
				FileValues:    params.FileValues,
				RootPath:      params.RootPath,
				SecretLoaders: params.SecretLoaders,
				LeftDelims:    params.LeftDelims,
				RightDelims:   params.RightDelims,
				AltDelims:     params.AltDelims,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.InputPath, "in", "-", "Template input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.OutputPath, "out", "", "Output file ('-' for stdout or a filename)")
	cmd.Flags().StringVar(&params.RootPath, "root", "", "Defines file loader root base path")
	cmd.Flags().StringArrayVarP(&params.SecretLoaders, "secrets-from", "s", []string{"vault"}, "Specifies secret containers to load ('vault' for Vault loader or '-' for stdin or filename)")
	cmd.Flags().StringArrayVarP(&params.ValueFiles, "values", "f", []string{}, "Specifies value files to load")
	cmd.Flags().StringArrayVar(&params.Values, "set", []string{}, "Specifies value (k=v)")
	cmd.Flags().StringArrayVar(&params.StringValues, "set-string", []string{}, "Specifies value (k=string)")
	cmd.Flags().StringArrayVar(&params.FileValues, "set-file", []string{}, "Specifies value (k=filepath)")
	cmd.Flags().StringVar(&params.LeftDelims, "left-delimiter", "{{", "Template left delimiter (default to '{{')")
	cmd.Flags().StringVar(&params.RightDelims, "right-delimiter", "}}", "Template right delimiter (default to '}}')")
	cmd.Flags().BoolVar(&params.AltDelims, "alt-delims", false, "Define '[[' and ']]' as template delimiters.")

	return cmd
}
