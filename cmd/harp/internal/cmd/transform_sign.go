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
	"zntr.io/harp/v2/pkg/sdk/value/signature"
)

// -----------------------------------------------------------------------------

type transformSignParams struct {
	inputPath   string
	outputPath  string
	keyRaw      string
	preHashed   bool
	detached    bool
	determistic bool
}

var transformSignCmd = func() *cobra.Command {
	params := &transformSignParams{}

	cmd := &cobra.Command{
		Use:     "sign",
		Short:   "Sign the given value with a transformer",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-sign", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Resolve tranformer
			t, err := signature.FromKey(params.keyRaw)
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

			// Transformation flag
			ctx = signature.WithDetachedSignature(ctx, params.detached)
			ctx = signature.WithDetermisticSignature(ctx, params.determistic)
			ctx = signature.WithInputPreHashed(ctx, params.preHashed)

			// Apply transformation
			out, err := t.To(ctx, content)
			if err != nil {
				log.For(ctx).Fatal("unable to apply transformer", zap.Error(err))
			}

			// Dump as output
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
	cmd.Flags().BoolVar(&params.detached, "detached", false, "Returns the signature only")
	cmd.Flags().BoolVar(&params.preHashed, "pre-hashed", false, "The input is already pre-hashed")
	cmd.Flags().BoolVar(&params.determistic, "deterministic", false, "Use deterministic signature algorithm variant (if key permits)")

	return cmd
}
