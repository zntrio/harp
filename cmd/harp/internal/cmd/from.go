// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var fromCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "from",
		Short: "Secret container generation commands",
	}

	// Add subcommands
	cmd.AddCommand(fromVaultCmd())
	cmd.AddCommand(fromJSONCmd())
	cmd.AddCommand(fromTemplateCmd())
	cmd.AddCommand(fromDumpCmd())
	cmd.AddCommand(fromOPLogCmd())
	cmd.AddCommand(fromObjectCmd())
	cmd.AddCommand(fromConsulCmd())
	cmd.AddCommand(fromEtcd3Cmd())
	cmd.AddCommand(fromZookeeperCmd())
	cmd.AddCommand(fromHCLCmd())

	return cmd
}
