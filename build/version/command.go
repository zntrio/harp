// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package version

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var (
	displayAsJSON bool
	withModules   bool
)

// Command exports Cobra command builder.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display service version",
		Run: func(cmd *cobra.Command, args []string) {
			bi := NewInfo()
			if displayAsJSON {
				fmt.Fprintf(os.Stdout, "%s", bi.JSON())
			} else {
				fmt.Fprintf(os.Stdout, "%s", bi.String())
				if withModules {
					fmt.Fprintln(os.Stdout, "\nDependencies:")
					for _, dep := range bi.BuildDeps {
						fmt.Fprintf(os.Stdout, "- %s\n", dep)
					}
				}
			}
		},
	}

	// Register parameters
	cmd.Flags().BoolVar(&displayAsJSON, "json", false, "Display build info as json")
	cmd.Flags().BoolVar(&withModules, "with-modules", false, "Display builtin go modules")

	// Return command
	return cmd
}
