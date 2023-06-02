// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/security/diceware"
)

var passphraseWordCount int8

// -----------------------------------------------------------------------------

var passphraseCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passphrase",
		Short: "Generate and print a diceware passphrase",
		Run:   runPassphrase,
	}

	// Parameters
	cmd.Flags().Int8VarP(&passphraseWordCount, "word-count", "w", 8, "Word count in diceware passphrase")

	return cmd
}

func runPassphrase(cmd *cobra.Command, args []string) {
	ctx, cancel := cmdutil.Context(cmd.Context(), "harp-passphrase", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// Check lower limit
	if passphraseWordCount < 4 {
		passphraseWordCount = 4
	}

	// Generate passphrase
	passPhrase, err := diceware.Diceware(int(passphraseWordCount))
	if err != nil {
		log.For(ctx).Fatal("unable to generate diceware passphrase", zap.Error(err))
	}

	// Print the key
	// lgtm [go/clear-text-logging]
	fmt.Fprintln(os.Stdout, passPhrase)
}
