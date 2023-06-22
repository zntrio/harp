// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// -----------------------------------------------------------------------------

func bugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bug",
		Short: "start a bug report",
		Long: `
	Bug opens the default browser and starts a new bug report.
	The report includes useful system information.
		`,
		Run: runBug,
	}
}

func runBug(cmd *cobra.Command, args []string) {
	ctx, cancel := cmdutil.Context(cmd.Context(), "harp-bug", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// No argument check
	if len(args) > 0 {
		log.For(ctx).Fatal("bug command takes no arguments")
	}

	// Prepare the report body
	body := cmdutil.BugReport()

	// Open the browser to issue creation form
	reportURL := "https://zntr.io/harp/issues/new?body=" + url.QueryEscape(body)
	if err := open.Run(reportURL); err != nil {
		fmt.Fprint(os.Stdout, "Please file a new issue at zntr.io/harp/issues/new using this template:\n\n")
		fmt.Fprint(os.Stdout, body)
	}
}

// -----------------------------------------------------------------------------
