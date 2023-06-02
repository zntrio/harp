// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var pluginCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage harp plugins",
	}

	// Add commands
	cmd.AddCommand(pluginListCmd())

	return cmd
}
