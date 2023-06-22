// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"zntr.io/harp/v2/build/fips"
)

// -----------------------------------------------------------------------------

var keygenCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keygen",
		Aliases: []string{"kg"},
		Short:   "Key generation commands",
	}

	// Subcommands
	cmd.AddCommand(keygenFernetCmd())
	cmd.AddCommand(keygenAESCmd())
	cmd.AddCommand(keygenMasterKeyCmd())
	cmd.AddCommand(keygenKeypairCmd())
	cmd.AddCommand(keygenPreSharedKeyCmd())

	if !fips.Enabled() {
		cmd.AddCommand(keygenSecretBoxCmd())
		cmd.AddCommand(keygenChaChaCmd())
		cmd.AddCommand(keygenXChaChaCmd())
		cmd.AddCommand(keygenAESPMACSIVCmd())
		cmd.AddCommand(keygenAESSIVCmd())
		cmd.AddCommand(keygenPasetoCmd())
		cmd.AddCommand(keygenBrancaCmd())
	}
	return cmd
}
