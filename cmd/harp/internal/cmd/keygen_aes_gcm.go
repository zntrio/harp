// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/awnumar/memguard"
	"github.com/spf13/cobra"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// -----------------------------------------------------------------------------

var keygenAESCmd = func() *cobra.Command {
	var keySize uint16

	cmd := &cobra.Command{
		Use:     "aes-gcm",
		Aliases: []string{"aes"},
		Short:   "Generate and print an AES-GCM key",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-aes", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Validate key size
			switch keySize {
			case 128, 192, 256:
				break
			default:
				log.For(ctx).Fatal("invalid specificed key size, only 128, 192 and 256 are supported.")
			}

			fmt.Fprintf(os.Stdout, "aes-gcm:%s", base64.URLEncoding.EncodeToString(memguard.NewBufferRandom(int(keySize/8)).Bytes()))
		},
	}

	// Parameters
	cmd.Flags().Uint16Var(&keySize, "size", 128, "Specify an AES key size (128, 192, 256)")

	return cmd
}
