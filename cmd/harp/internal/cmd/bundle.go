// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var bundleCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bundle",
		Aliases: []string{"b"},
		Short:   "Bundle commands",
	}

	// Bundle commands
	cmd.AddCommand(bundleDumpCmd())
	cmd.AddCommand(bundleReadCmd())
	cmd.AddCommand(bundleEncryptCmd())
	cmd.AddCommand(bundleDecryptCmd())
	cmd.AddCommand(bundleDiffCmd())
	cmd.AddCommand(bundlePatchCmd())
	cmd.AddCommand(bundleFilterCmd())
	cmd.AddCommand(bundleLintCmd())
	cmd.AddCommand(bundlePrefixerCmd())

	return cmd
}
