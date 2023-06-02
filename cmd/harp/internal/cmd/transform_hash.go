// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/security"
	"zntr.io/harp/v2/pkg/sdk/value/encoding"
	"zntr.io/harp/v2/pkg/sdk/value/hash"
)

// -----------------------------------------------------------------------------

type transformHashParams struct {
	inputPath  string
	outputPath string
	algorithm  string
	encoding   string
	validate   string
}

var transformHashCmd = func() *cobra.Command {
	params := &transformHashParams{}

	longDesc := cmdutil.LongDesc(fmt.Sprintf(`
		Process the input to compute the hash according to selected hash algorithm.

		The command input is limited to size lower than 250 MB.

		Supported Algorithms:
		  %s`, strings.Join(hash.SupportedAlgorithms(), ", ")))

	examples := cmdutil.Examples(`
		# Compute SHA256 from stdin
		echo -n 'test' | harp transform hash

		# Compute SHA512 hash from a file
		harp transform hash --algorithm sha512

		# Compute Blake2b hash from a file with base64 encoded output
		harp transform hash --algorithm blake2b-512 --encoding base64

		# Check the given input integrity (default sha256 / hex)
		harp transform hash --in livecd.iso --validate 4506369c20d2a95ebad9234b7f48e0eded4ec4ee1de0cb45a195b1e38fde27f7

		# Check the given input integrity with specific hash algorihm and encoding
		harp transform hash --in livecd.iso --algorithm BLAKE2b_512 --encoding base64urlraw --validate dquOtQ-gj815njSbk8mGl3WUgImkflX1AaLXy6ymhk_kUpP6qXDmSC5X2l3nkTgJK9F6p3rBV6o075QZQ-HHaw`)

	cmd := &cobra.Command{
		Use:     "hash",
		Short:   "Hash given input",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-hash", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
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

			// Prepare hasher
			h, err := hash.NewHasher(params.algorithm)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize hasher", zap.Error(err))
			}

			// Read chunk
			if _, err := io.Copy(h, reader); err != nil {
				log.For(ctx).Fatal("unable to read content chunk", zap.Error(err))
			}

			// Finalize
			content := h.Sum(nil)

			// Validation mode
			if params.validate != "" {
				encoderReader, err := encoding.NewReader(strings.NewReader(params.validate), params.encoding)
				if err != nil {
					log.For(ctx).Fatal("unable to decode expected hash", zap.Error(err))
				}

				expectedHash, err := io.ReadAll(encoderReader)
				if err != nil {
					log.For(ctx).Fatal("unable to read expected hash", zap.Error(err))
				}

				// Compare expected and given
				if !security.SecureCompare(expectedHash, content) {
					log.For(ctx).Fatal("invalid expectation", zap.String("expected", params.validate), zap.String("got", hex.EncodeToString(content)))
				}

				// Dump as output
				if _, err := fmt.Fprintln(writer, "OK"); err != nil {
					log.For(ctx).Fatal("unable to write comparison result", zap.Error(err))
				}
			} else {
				encoderWriter, err := encoding.NewWriter(writer, params.encoding)
				if err != nil {
					log.For(ctx).Fatal("unable to write content hash", zap.Error(err))
				}

				// Process input as a stream.
				if _, err := io.Copy(encoderWriter, bytes.NewReader(content)); err != nil {
					log.SafeClose(encoderWriter, "unable to close the encoder writer")
					log.For(ctx).Fatal("unable to process input", zap.Error(err))
				}

				log.SafeClose(encoderWriter, "unable to close the encoder writer")
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.algorithm, "algorithm", "SHA256", "Hash algorithm to use")
	cmd.Flags().StringVar(&params.encoding, "encoding", "hex", "Encoding strategy (hex, base64, base64raw, base64url, base64urlraw)")
	cmd.Flags().StringVar(&params.validate, "validate", "", "Expecetd hash to validate the output with. Decoded using the given encoding strategy.")

	return cmd
}
