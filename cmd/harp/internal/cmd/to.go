// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var toCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to",
		Short: "Secret container conversion commands",
	}

	// Add sub commands
	cmd.AddCommand(toVaultCmd())
	cmd.AddCommand(toObjectCmd())
	cmd.AddCommand(toRulesetCmd())
	cmd.AddCommand(toEtcd3Cmd())
	cmd.AddCommand(toConsulCmd())
	cmd.AddCommand(toZookeeperCmd())
	cmd.AddCommand(toGithubActionCmd())

	return cmd
}
