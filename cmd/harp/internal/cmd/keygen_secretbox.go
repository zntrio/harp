// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/awnumar/memguard"
	"github.com/spf13/cobra"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
)

// -----------------------------------------------------------------------------

var keygenSecretBoxCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secretbox",
		Aliases: []string{"sb"},
		Short:   "Generate and print a secretbox keypair",
		Run:     runKeygenSecretbox,
	}

	return cmd
}

func runKeygenSecretbox(cmd *cobra.Command, args []string) {
	_, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-secretbox", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	fmt.Fprintf(os.Stdout, "secretbox:%s", base64.URLEncoding.EncodeToString(memguard.NewBufferRandom(32).Bytes()))
}
