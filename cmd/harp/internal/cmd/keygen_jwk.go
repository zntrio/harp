// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/go-jose/go-jose/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/keygen"
)

// -----------------------------------------------------------------------------.
type keygenJWKParams struct {
	outputPath         string
	signatureAlgorithm string
	keyBits            int
	keyID              string
}

var keygenKeypairCmd = func() *cobra.Command {
	params := &keygenJWKParams{}

	cmd := &cobra.Command{
		Use:   "jwk",
		Short: "Generate a JWK encoded key pair",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-jwk", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &keygen.JWKTask{
				OutputWriter:       cmdutil.FileWriter(params.outputPath),
				SignatureAlgorithm: params.signatureAlgorithm,
				KeySize:            params.keyBits,
				KeyID:              params.keyID,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Add parameters
	cmd.Flags().StringVar(&params.signatureAlgorithm, "algorithm", string(jose.EdDSA), "Key type to generate")
	cmd.Flags().IntVar(&params.keyBits, "bits", 0, "Key size (in bits)")
	cmd.Flags().StringVar(&params.keyID, "key-id", "", "Key identifier")

	return cmd
}
