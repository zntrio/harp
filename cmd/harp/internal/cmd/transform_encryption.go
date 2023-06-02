// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

// -----------------------------------------------------------------------------

var transformEncryptionCmd = func() *cobra.Command {
	var (
		inputPath  string
		outputPath string
		keyRaw     string
		revert     bool
	)

	cmd := &cobra.Command{
		Use:        "encryption",
		Short:      "Encryption value transformer",
		Aliases:    []string{"enc"},
		Deprecated: "Use encrypt/decrypt commands.",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-encryption", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Resolve tranformer
			t, err := encryption.FromKey(keyRaw)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize a transformer form key", zap.Error(err))
			}
			if t == nil {
				log.For(ctx).Fatal("transformer is nil")
			}

			// Read input
			reader, err := cmdutil.Reader(inputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize input reader", zap.Error(err))
			}

			// Read input
			writer, err := cmdutil.Writer(inputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize output writer", zap.Error(err))
			}

			// Drain reader
			content, err := io.ReadAll(reader)
			if err != nil {
				log.For(ctx).Fatal("unable to drain input reader", zap.Error(err))
			}

			var out []byte
			if !revert {
				// Apply transformation
				out, err = t.To(ctx, content)
				if err != nil {
					log.For(ctx).Fatal("unable to apply transformer", zap.Error(err))
				}
			} else {
				// Apply transformation
				out, err = t.From(ctx, content)
				if err != nil {
					log.For(ctx).Fatal("unable to apply transformer", zap.Error(err))
				}
			}

			if _, err = writer.Write(out); err != nil {
				log.For(ctx).Fatal("unable to write result to writer", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&keyRaw, "key", "", "Transformer key")
	log.CheckErr("unable to mark 'key' flag as required.", cmd.MarkFlagRequired("key"))

	cmd.Flags().StringVar(&inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().BoolVar(&revert, "revert", false, "Decrypt the input (default encrypt)")

	return cmd
}
