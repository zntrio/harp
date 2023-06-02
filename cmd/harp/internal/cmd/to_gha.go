// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/tasks/to"
)

// -----------------------------------------------------------------------------

type toGithubActionParams struct {
	inputPath    string
	owner        string
	repository   string
	secretFilter string
}

var toGithubActionCmd = func() *cobra.Command {
	var params toGithubActionParams

	cmd := &cobra.Command{
		Use:     "github-actions",
		Aliases: []string{"gha"},
		Short:   "Export all secrets to Github Actions as repository secrets.",
		Example: `$ export GITHUB_TOKEN=ghp_###############
$ harp to gha --in secret.container --owner elastic --owner harp --secret-filter "COSIGN_*"`,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-to-gha", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			// Prepare task
			t := &to.GithubActionTask{
				ContainerReader: cmdutil.FileReader(params.inputPath),
				Owner:           params.owner,
				Repository:      params.repository,
				SecretFilter:    params.secretFilter,
			}

			// Run the task
			if err := t.Run(ctx); err != nil {
				log.For(ctx).Fatal("unable to execute task", zap.Error(err))
			}
		},
	}

	// Parameters
	cmd.Flags().StringVar(&params.inputPath, "in", "-", "Container path ('-' for stdin or filename)")
	cmd.Flags().StringVar(&params.owner, "owner", "", "Github owner/organization")
	cmd.Flags().StringVar(&params.repository, "repository", "", "Github repository")
	cmd.Flags().StringVar(&params.secretFilter, "secret-filter", "*", "Specify secret filter as Glob (*_KEY, private*)")

	return cmd
}
