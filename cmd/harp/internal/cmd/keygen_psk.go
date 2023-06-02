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

var keygenPreSharedKeyCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pre-shared-key",
		Aliases: []string{"psk"},
		Short:   "Generate and print a container pre-shared-key",
		Run:     runKeygenPreSharedKey,
	}

	return cmd
}

func runKeygenPreSharedKey(cmd *cobra.Command, args []string) {
	_, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-psk", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	fmt.Fprintf(os.Stdout, "%s", base64.RawURLEncoding.EncodeToString(memguard.NewBufferRandom(64).Bytes()))
}
