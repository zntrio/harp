// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var csoCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cso",
		Short: "CSO commands",
	}

	// Sub-commands
	cmd.AddCommand(csoValidateCmd())
	cmd.AddCommand(csoParseCmd())

	return cmd
}
