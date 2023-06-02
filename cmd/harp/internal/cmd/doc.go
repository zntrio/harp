// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
)

// -----------------------------------------------------------------------------

var docCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Generates documentation and autocompletion",
	}

	// Subcommands
	cmd.AddCommand(docMarkdownCmd())

	return cmd
}

// -----------------------------------------------------------------------------

var docDestination string

var docMarkdownCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "markdown",
		Aliases: []string{"md"},
		Short:   "Documentation in Markdown format",
		RunE:    runDocMarkdown,
	}

	// Parameters
	cmd.Flags().StringVarP(&docDestination, "destination", "d", "", "destination for documentation")

	return cmd
}

func runDocMarkdown(cmd *cobra.Command, args []string) error {
	// Context to attach all goroutines
	_, cancel := cmdutil.Context(cmd.Context(), "harp-doc-markdown", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// Disable flag
	cmd.Root().DisableAutoGenTag = true

	// Generate markdown tree
	return doc.GenMarkdownTree(cmd.Root(), docDestination)
}
