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
	"zntr.io/harp/v2/build/version"
	iconfig "zntr.io/harp/v2/cmd/harp/internal/config"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/config"
	configcmd "zntr.io/harp/v2/pkg/sdk/config/cmd"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// -----------------------------------------------------------------------------

var (
	cfgFile string
	conf    = &iconfig.Configuration{}
)

// -----------------------------------------------------------------------------

// RootCmd describes root command of the tool.
var mainCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harp",
		Short: "Extensible secret management tool",
	}

	// Register flags
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// Register sub commands
	cmd.AddCommand(version.Command())
	cmd.AddCommand(configcmd.NewConfigCommand(conf, "HARP"))

	cmd.AddCommand(bundleCmd())
	cmd.AddCommand(containerCmd())
	cmd.AddCommand(keygenCmd())
	cmd.AddCommand(passphraseCmd())
	cmd.AddCommand(docCmd())
	cmd.AddCommand(bugCmd())

	cmd.AddCommand(pluginCmd())
	cmd.AddCommand(csoCmd())

	cmd.AddCommand(templateCmd())
	cmd.AddCommand(renderCmd())
	cmd.AddCommand(valuesCmd())

	cmd.AddCommand(fromCmd())
	cmd.AddCommand(toCmd())

	cmd.AddCommand(transformCmd())
	cmd.AddCommand(shareCmd())
	cmd.AddCommand(lintCmd())

	// Return command
	return cmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

// -----------------------------------------------------------------------------

// Execute main command.
func Execute() error {
	args := os.Args

	// Initialize root command
	cmd := mainCmd()

	// Initialize plugin handler
	pluginHandler := cmdutil.NewDefaultPluginHandler(validPluginFilenamePrefixes)

	// If has more than 1 arguments
	if len(args) > 1 {
		cmdPathPieces := args[1:]

		// only look for suitable extension executables if
		// the specified command does not already exist
		if _, _, err := cmd.Find(cmdPathPieces); err != nil {
			if err := cmdutil.HandlePluginCommand(pluginHandler, cmdPathPieces); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}
	}

	return cmd.Execute()
}

// -----------------------------------------------------------------------------

func initConfig() {
	if err := config.Load(conf, "HARP", cfgFile); err != nil {
		log.Bg().Fatal("Unable to load settings", zap.Error(err))
	}
}
