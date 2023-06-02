// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------

var transformCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Transform input value using encryption transformers",
	}

	// Deprecated
	cmd.AddCommand(transformEncryptionCmd())

	// Add commands
	cmd.AddCommand(transformEncryptCmd())
	cmd.AddCommand(transformDecryptCmd())
	cmd.AddCommand(transformSignCmd())
	cmd.AddCommand(transformVerifyCmd())
	cmd.AddCommand(transformDecodeCmd())
	cmd.AddCommand(transformEncodeCmd())
	cmd.AddCommand(transformHashCmd())
	cmd.AddCommand(transformCompressCmd())
	cmd.AddCommand(transformDecompressCmd())
	cmd.AddCommand(transformMultihashCmd())

	return cmd
}
