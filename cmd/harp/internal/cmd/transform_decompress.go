// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/ioutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/value/compression"
)

// -----------------------------------------------------------------------------

type transformDecompressParams struct {
	inputPath             string
	outputPath            string
	algorithm             string
	maxDecompressionGuard uint16
}

var transformDecompressCmd = func() *cobra.Command {
	params := &transformDecompressParams{}

	longDesc := cmdutil.LongDesc(`
	Decompress the given input stream using the selected compression algorithm.

	Supported compression:
	  * identity - returns the unmodified input
	  * gzip
	  * lzw/lzw-msb/lzw-lsb
	  * lz4
	  * s2/snappy
	  * zlib
	  * flate/deflate
	  * lzma
	  * zstd`)

	examples := cmdutil.Examples(`
	# Compress a file
	harp transform decompress --in README.md.gz --out README.md --algorithm gzip

	# Decompress to STDOUT
	harp transform compress --in README.md.gz --algorithm gzip

	# Decompress from STDIN
	harp transform compress --algorithm gzip`)

	cmd := &cobra.Command{
		Use:     "decompress",
		Short:   "Decompress given input",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-decompress", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
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

			// Prepare compressor
			compressedReader, err := compression.NewReader(reader, params.algorithm)
			if err != nil {
				log.SafeClose(compressedReader, "unable to close the compression writer")
				log.For(ctx).Fatal("unable to write encoded content", zap.Error(err))
			}

			// Compute max decompression size
			maxDecompressionSize := uint64(params.maxDecompressionGuard) * 1024 * 1024

			// Process input as a stream.
			if _, err := ioutil.LimitCopy(writer, compressedReader, maxDecompressionSize); err != nil {
				log.SafeClose(compressedReader, "unable to close the compression writer")
				log.For(ctx).Fatal("unable to process input", zap.Error(err))
			}

			// Close the writer
			log.SafeClose(compressedReader, "unable to close the compression writer")
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.algorithm, "algorithm", "gzip", "Compression algorithm")
	cmd.Flags().Uint16Var(&params.maxDecompressionGuard, "max-decompression-guard", 100, "Decompression guard in MB")

	return cmd
}
