// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/value/hash"
)

// -----------------------------------------------------------------------------

type transformMultihashParams struct {
	inputPath  string
	outputPath string
	algorithms []string
	jsonOutput bool
}

var transformMultihashCmd = func() *cobra.Command {
	params := &transformMultihashParams{}

	longDesc := cmdutil.LongDesc(fmt.Sprintf(`
		Process the input to compute the hashes according to selected hash algorithms.

		The command input is limited to size lower than 250 MB.

		Supported Algorithms:
		  %s`, strings.Join(hash.SupportedAlgorithms(), ", ")))

	examples := cmdutil.Examples(`
	# Compute md5, sha1, sha256, sha512 in one read from a file
	harp transform multihash --in livecd.iso

	# Compute sha256, sha512 only
	harp transform multihash --algorithm sha256 --algorithm sha512 --in livecd.iso

	# Compute sha256, sha512 only with JSON output
	harp transform multihash --json --algorithm sha256 --algorithm sha512 --in livecd.iso
	`)

	cmd := &cobra.Command{
		Use:     "multihash",
		Aliases: []string{"mh"},
		Short:   "Multiple hash  from given input",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-transform-multihash", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
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
			hMap, err := hash.NewMultiHash(reader, params.algorithms...)
			if err != nil {
				log.For(ctx).Fatal("unable to initialize hasher", zap.Error(err))
			}

			// Display as json
			if params.jsonOutput {
				if err := json.NewEncoder(writer).Encode(hMap); err != nil {
					log.For(ctx).Fatal("unable to encode result as json", zap.Error(err))
				}
			} else {
				// Sort map keys to get a stable output.
				keys := make([]string, 0, len(hMap))
				for k := range hMap {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				for _, k := range keys {
					// Display container key
					if _, err := fmt.Fprintf(writer, "%s=%s\n", k, hMap[k]); err != nil {
						log.For(ctx).Fatal("unable to display result", zap.Error(err))
					}
				}
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Input path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Output path ('-' for stdout or filename)")
	cmd.Flags().StringSliceVar(&params.algorithms, "algorithm", []string{"md5", "sha1", "sha256", "sha512"}, "Hash algorithms to use")
	cmd.Flags().BoolVar(&params.jsonOutput, "json", false, "Display multihash result as json")

	return cmd
}
