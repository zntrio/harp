// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/from"
)

// -----------------------------------------------------------------------------

var fromHCLCmd = func() *cobra.Command {
	var (
		inputPath  string
		outputPath string
	)
	cmd := &cobra.Command{
		Use:   "hcl",
		Short: "Convert a HCL based bundle DSL to a secret-container",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-from-hcl", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &from.HCLTask{
				HCLReader:    cmdutil.FileReader(inputPath),
				OutputWriter: cmdutil.FileWriter(outputPath),
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&inputPath, "in", "-", "HCL configuration ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "-", "Container output ('-' for stdout or filename)")

	return cmd
}
