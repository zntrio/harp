// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/value/encoding"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
)

// -----------------------------------------------------------------------------

type transformEncryptParams struct {
	inputPath              string
	outputPath             string
	keyRaw                 string
	additionalData         string
	additionalDataEncoding string
}

var transformEncryptCmd = func() *cobra.Command {
	params := &transformEncryptParams{}

	cmd := &cobra.Command{
		Use:     "encrypt",
		Short:   "Encrypt the given value with a value transformer",
		Aliases: []string{"e"},
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-decrypt", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Resolve tranformer
			t, err := encryption.FromKey(params.keyRaw)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize a transformer form key", zap.Error(err))
			}
			if t == nil {
				log.For(ctx).Fatal("transformer is nil")
			}

			// Read input
			reader, err := cmdutil.Reader(params.inputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize input reader", zap.Error(err))
			}

			// Read input
			writer, err := cmdutil.Writer(params.outputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize output writer", zap.Error(err))
			}

			// Drain reader
			content, err := io.ReadAll(reader)
			if err != nil {
				log.For(ctx).Fatal("unable to drain input reader", zap.Error(err))
			}

			// Decode AAD if any
			if params.additionalData != "" {
				encoderReader, errDecode := encoding.NewReader(strings.NewReader(params.additionalData), params.additionalDataEncoding)
				if errDecode != nil {
					log.For(ctx).Fatal("unable to decode additional data", zap.Error(errDecode))
				}

				aad, errAADRead := io.ReadAll(encoderReader)
				if errAADRead != nil {
					log.For(ctx).Fatal("unable to read additional data", zap.Error(errAADRead))
				}

				// Set additional data
				ctx = encryption.WithAdditionalData(ctx, aad)
			}

			// Apply transformation
			out, err := t.To(ctx, content)
			if err != nil {
				log.For(ctx).Fatal("unable to apply transformer", zap.Error(err))
			}

			if _, err = writer.Write(out); err != nil {
				log.For(ctx).Fatal("unable to write result to writer", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.keyRaw, "key", "", "Transformer key")
	log.CheckErr("unable to mark 'key' flag as required.", cmd.MarkFlagRequired("key"))

	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.additionalData, "aad", "", "Additional data for AEAD encryption")
	cmd.Flags().StringVar(&params.additionalDataEncoding, "aad-encoding", "base64", "Additional data encoding strategy")

	return cmd
}
