// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/awnumar/memguard"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/container"
)

// -----------------------------------------------------------------------------.
type containerUnsealParams struct {
	inputPath       string
	outputPath      string
	containerKeyRaw string
	preSharedKeyRaw string
}

var containerUnsealCmd = func() *cobra.Command {
	params := containerUnsealParams{}

	cmd := &cobra.Command{
		Use:   "unseal",
		Short: "Unseal a secret container",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-container-unseal", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare passphrase
			containerKey := memguard.NewBufferFromBytes([]byte(params.containerKeyRaw))
			if params.containerKeyRaw == "" {
				var err error
				// Read passphrase from stdin
				containerKey, err = cmdutil.ReadSecret("Enter container key", false)
				if err != nil {
					log.For(ctx).Fatal("unable to read passphrase", zap.Error(err))
				}
			}
			defer containerKey.Destroy()

			// Prepare task
			t := &container.UnsealTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				OutputWriter:    cmdutil.StdoutWriter(),
				ContainerKey:    containerKey,
			}
			if params.preSharedKeyRaw != "" {
				t.PreSharedKey = memguard.NewBufferFromBytes([]byte(params.preSharedKeyRaw))
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "", "Sealed container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "", "Unsealed container output ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.containerKeyRaw, "key", "", "Container key")
	log.CheckErr("unable to mark 'key' flag as required.", cmd.MarkFlagRequired("key"))
	cmd.Flags().StringVar(&params.preSharedKeyRaw, "pre-shared-key", "", "Use a pre-shared-key to unseal the container")

	return cmd
}
