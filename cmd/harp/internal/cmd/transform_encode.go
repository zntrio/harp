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

type transformEncodeParams struct {
	inputPath  string
	outputPath string
	encoding   string
}

var transformEncodeCmd = func() *cobra.Command {
	params := &transformEncodeParams{}

	longDesc := cmdutil.LongDesc(`
	Encode the given input stream using the selected encoding strategy.

	Supported codecs:
	  * identity - returns the unmodified input
	  * hex/base16 - returns the hexadecimal encoded input
	  * base32 - returns the Base32 encoded input
	  * base32hex - returns the Base32 with extended alphabet encoded input
	  * base62 - returns the Base62 encoded input
	  * base64 - returns the Base64 encoded input
	  * base64raw - returns the Base64 encoded input without "=" padding
	  * base64url - returns the Base64 encoded input using URL safe characters
	  * base64urlraw - returns the Base64 encoded input using URL safe characters without "=" padding
	  * base85 - returns the Base85 encoded input`)

	examples := cmdutil.Examples(`
		# Encode stdin to base64
		echo "test" | harp transform encode --encoding base64

		# Encode a file using base64url
		harp transform encode --encoding base64url --in test.txt --out encoded.text `)

	cmd := &cobra.Command{
		Use:     "encode",
		Short:   "Encode given input",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-encode", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

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

			// Write encoded content
			encoderWriter, err := encoding.NewWriter(writer, params.encoding)
			if err != nil {
				log.SafeClose(encoderWriter, "unable to close the encoder writer")
				log.For(ctx).Fatal("unable to write encoded content", zap.Error(err))
			}

			// Process input as a stream.
			if _, err := io.Copy(encoderWriter, reader); err != nil {
				log.SafeClose(encoderWriter, "unable to close the encoder writer")
				log.For(ctx).Fatal("unable to process input", zap.Error(err))
			}

			// Close the writer
			log.SafeClose(encoderWriter, "unable to close the encoder writer")
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.encoding, "encoding", "identity", "Encoding strategy")

	return cmd
}
