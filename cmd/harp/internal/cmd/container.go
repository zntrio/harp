// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var containerCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "container",
		Aliases: []string{"c"},
		Short:   "Secret container commands",
	}

	// Bundle commands
	cmd.AddCommand(containerIdentityCmd())
	cmd.AddCommand(containerRecoveryCmd())
	cmd.AddCommand(containerSealCmd())
	cmd.AddCommand(containerUnsealCmd())

	return cmd
}
