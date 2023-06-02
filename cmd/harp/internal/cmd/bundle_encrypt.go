// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
	"zntr.io/harp/v2/pkg/tasks/bundle"
)

// -----------------------------------------------------------------------------.
type bundleEncryptParams struct {
	inputPath      string
	outputPath     string
	key            string
	keyAliases     []string
	skipUnresolved bool
}

var bundleEncryptCmd = func() *cobra.Command {
	params := &bundleEncryptParams{}

	longDesc := cmdutil.LongDesc(`
	Apply package content encryption.

	For confidentiality purpose, bundle package value can be encrypted before
	the container sealing. It offers confidentiality properties so that the
	final consumer must know an additional decryption key to be allowed to
	read the package value even if it can unseal the container.

	All package properties (name, labels, annotations) remain a clear-text
	message. Only package values (secret K/V) are encrypted.

	This act as in-transit/in-use encryption.

	Annotations:

	* harp.zntr.io/v2/package#encryptionKeyAlias=<alias> - Set this
	  annotation on packages to reference a key alias.`)

	examples := cmdutil.Examples(`
	# Encrypt a whole bundle from STDIN and produce output to STDOUT
	harp bundle encrypt --key <transformer key>

	# Encrypt partially a bundle using the annotation matcher from STDIN and
	# produce output to STDOUT
	harp bundle encrypt --key-alias <alias>:<transformer key> --key-alias <alias-2>:<transformer key 2>
	`)

	cmd := &cobra.Command{
		Use:     "encrypt",
		Short:   "Encrypt secret values",
		Long:    longDesc,
		Example: examples,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bundle-encrypt", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &bundle.EncryptTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				OutputWriter:    cmdutil.FileWriter(params.outputPath),
			}
			switch {
			case params.key != "":
				// Create transformer according to used encryption key
				transformer, err := encryption.FromKey(params.key)
				if err != nil {
					log.For(ctx).Fatal("unable to initialize transformer", zap.Error(err))
				}

				// Use the given key a bundle transformer
				t.BundleTransformer = transformer
			case len(params.keyAliases) > 0:
				transformerMap := map[string]value.Transformer{}

				// Split all alias / key
				for _, alias := range params.keyAliases {
					// Split alias
					parts := strings.SplitN(alias, ":", 2)
					if len(parts) != 2 {
						log.For(ctx).Fatal("invalid alias, it must be formatted alias:key.", zap.String("alias", alias))
						return
					}

					// Create transformer according to used encryption key
					transformer, err := encryption.FromKey(parts[1])
					if err != nil {
						log.For(ctx).Fatal("unable to initialize transformer", zap.Error(err))
					}

					// Assign to map
					transformerMap[parts[0]] = transformer
				}

				// Use transformer map
				t.TransformerMap = transformerMap
				t.SkipUnresolved = params.skipUnresolved
			default:
				log.For(ctx).Fatal("--key or --key-alias must be provided")
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "", "Container input ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.outputPath, "out", "", "Container output ('-' for stdout or filename)")
	cmd.Flags().StringVar(&params.key, "key", "", "Secret value encryption key for full bundle encryption")
	cmd.Flags().StringSliceVar(&params.keyAliases, "key-alias", []string{}, "Secret value encryption key for partial bundle encryption ('alias:key')")
	cmd.Flags().BoolVarP(&params.skipUnresolved, "skip-unresolved-key-alias", "s", false, "Skip unresolved key alias during partial bundle encryption")

	return cmd
}
