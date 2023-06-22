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
	"zntr.io/harp/v2/pkg/sdk/value/encoding"
)

// -----------------------------------------------------------------------------

type transformDecodeParams struct {
	inputPath  string
	outputPath string
	encoding   string
}

var transformDecodeCmd = func() *cobra.Command {
	params := &transformDecodeParams{}

	longDesc := cmdutil.LongDesc(`
	Decode the given input stream using the selected decoding strategy.

	Supported codecs:
	  * identity - returns the unmodified input
	  * hex/base16 - returns the hexadecimal decoded input
	  * base32 - returns the Base32 decoded input
	  * base32hex - returns the Base32 with extended alphabet decoded input
	  * base62 - returns the Base62 decoded input
	  * base64 - returns the Base64 decoded input
	  * base64raw - returns the Base64 decoded input without "=" padding
	  * base64url - returns the Base64 decoded input using URL safe characters
	  * base64urlraw - returns the Base64 decoded input using URL safe characters without "=" padding
	  * base85 - returns the Base85 decoded input`)

	examples := cmdutil.Examples(`
		# Decode base64 from stdin
		echo "dGVzdAo=" | harp transform decode --encoding base64

		# Decode base64url from a file
		harp transform decode --in test.txt --encoding base64url`)

	cmd := &cobra.Command{
		Use:     "decode",
		Short:   "Decode given input",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-decode", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Read input
			reader, err := cmdutil.Reader(params.inputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize input reader", zap.Error(err))
			}

			// Output writer
			writer, err := cmdutil.Writer(params.outputPath)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize output writer", zap.Error(err))
			}

			// Read and decode
			out, err := encoding.NewReader(reader, params.encoding)
			if err != nil {
				log.For(ctx).Fatal("unable to prepare input decoder", zap.Error(err))
			}

			// Process input as a stream.
			if _, err := io.Copy(writer, out); err != nil {
				log.For(ctx).Fatal("unable to process input", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.encoding, "encoding", "identity", "Encoding strategy")

	return cmd
}
