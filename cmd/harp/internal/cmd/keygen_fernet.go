// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/fernet/fernet-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// -----------------------------------------------------------------------------

var keygenFernetCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fernet",
		Short: "Generate and print a fernet key",
		Run:   runKeygenFernet,
	}

	return cmd
}

func runKeygenFernet(cmd *cobra.Command, args []string) {
	ctx, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-fernet", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// Generate a fernet key
	k := &fernet.Key{}
	if err := k.Generate(); err != nil {
		log.For(ctx).Fatal("unable to generate Fernet key", zap.Error(err))
	}

	// Print the key
	fmt.Fprint(os.Stdout, k.Encode())
}
