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

var fromVaultCmd = func() *cobra.Command {
	var (
		pathsFrom         string
		secretPaths       []string
		outputPath        string
		namespace         string
		withMetadata      bool
		withVaultMetadata bool
		maxWorkerCount    int64
		continueOnError   bool
	)

	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Pull a list of Vault K/V paths as a secret container",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-from-vault", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Check if we have to read external path
			if pathsFrom != "" {
				// Force read from stdin
				paths, errReader := cmdutil.LineReader(pathsFrom)
				if errReader != nil {
					log.For(ctx).Fatal("unable to read paths from stdin", zap.Error(errReader))
				}

				// Add to paths
				secretPaths = append(secretPaths, paths...)
			}

			// Prepare task
			t := &from.VaultTask{
				OutputWriter:    cmdutil.FileWriter(outputPath),
				SecretPaths:     secretPaths,
				VaultNamespace:  namespace,
				WithMetadata:    withMetadata || withVaultMetadata,
				AsVaultMetadata: withVaultMetadata,
				MaxWorkerCount:  maxWorkerCount,
				ContinueOnError: continueOnError,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&pathsFrom, "paths-from", "", "Path to read path from ('-' for stdin or filename)")
	cmd.Flags().StringArrayVar(&secretPaths, "path", []string{}, "Vault backend path (and recursive)")
	cmd.Flags().StringVar(&outputPath, "out", "", "Container output ('-' for stdout or filename)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Vault namespace")
	cmd.Flags().BoolVar(&withMetadata, "with-metadata", false, "Push container metadata as secret data")
	cmd.Flags().BoolVar(&withVaultMetadata, "with-vault-metadata", false, "Push container metadata as secret metadata (requires Vault >=1.9)")
	cmd.Flags().Int64Var(&maxWorkerCount, "worker-count", 4, "Active worker count limit")
	cmd.Flags().BoolVar(&continueOnError, "continue-on-error", false, "Continue exploration even when there is raised errors (permission denied)")

	return cmd
}
