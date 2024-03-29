// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/from"
	tplcmdutil "zntr.io/harp/v2/pkg/template/cmdutil"
	"zntr.io/harp/v2/pkg/template/engine"
)

// -----------------------------------------------------------------------------

var fromTemplateCmd = func() *cobra.Command {
	var (
		inputPath    string
		outputPath   string
		rootPath     string
		valueFiles   []string
		values       []string
		stringValues []string
		fileValues   []string
	)

	cmd := &cobra.Command{
		Use:     "bundle-template",
		Aliases: []string{"template", "t", "tmpl"},
		Short:   "Generate a secret container from a BundleTemplate manifest",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-from-template", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Load values
			valueOpts := tplcmdutil.ValueOptions{
				ValueFiles:   valueFiles,
				Values:       values,
				StringValues: stringValues,
				FileValues:   fileValues,
			}
			values, err := valueOpts.MergeValues()
			if err != nil {
				log.For(ctx).Fatal("unable to process values", zap.Error(err))
			}

			// Load files
			var files engine.Files
			if rootPath != "" {
				absRootPath, err := filepath.Abs(rootPath)
				if err != nil {
					log.For(ctx).Fatal("unable to get absolute file paht for root path", zap.Error(err))
				}

				files, err = tplcmdutil.Files(os.DirFS(absRootPath), ".")
				if err != nil {
					log.For(ctx).Fatal("unable to process files", zap.Error(err))
				}
			}

			// Prepare task
			t := &from.BundleTemplateTask{
				TemplateReader: cmdutil.FileReader(inputPath),
				OutputWriter:   cmdutil.FileWriter(outputPath),
				TemplateContext: engine.NewContext(
					engine.WithName(inputPath),
					engine.WithValues(values),
					engine.WithFiles(files),
				),
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "Template input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "", "Container output ('-' for stdout or a filename)")
	cmd.Flags().StringVar(&rootPath, "root", "", "Defines file loader root base path")
	cmd.Flags().StringArrayVarP(&valueFiles, "values", "f", []string{}, "Specifies value files to load")
	cmd.Flags().StringArrayVar(&values, "set", []string{}, "Specifies value (k=v)")
	cmd.Flags().StringArrayVar(&stringValues, "set-string", []string{}, "Specifies value (k=string)")
	cmd.Flags().StringArrayVar(&fileValues, "set-file", []string{}, "Specifies value (k=filepath)")

	return cmd
}
