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

var keygenAESPMACSIVCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aes-pmac-siv",
		Short: "Generate and print an AES-PMAC-SIV key",
		Run: func(cmd *cobra.Command, args []string) {
			_, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-aespmacsiv", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			fmt.Fprintf(os.Stdout, "aes-pmac-siv:%s", base64.URLEncoding.EncodeToString(memguard.NewBufferRandom(64).Bytes()))
		},
	}

	return cmd
}
