// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var shareCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Share secret using Vault Cubbyhole",
	}

	// Add sub commands
	cmd.AddCommand(sharePutCmd())
	cmd.AddCommand(shareGetCmd())

	return cmd
}
