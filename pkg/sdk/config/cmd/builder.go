// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/sdk/flags"
	"zntr.io/harp/v2/pkg/sdk/log"
)

var configNewAsEnvFlag bool

// NewConfigCommand initialize a cobra config command tree.
func NewConfigCommand(conf interface{}, envPrefix string) *cobra.Command {
	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// config
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Service Configuration",
	}

	// config new
	configNewCmd := &cobra.Command{
		Use:   "new",
		Short: "Initialize a default configuration",
		Run: func(cmd *cobra.Command, args []string) {
			defaults.SetDefaults(conf)

			if !configNewAsEnvFlag {
				btes, err := toml.Marshal(conf)
				if err != nil {
					log.For(cmd.Context()).Fatal("Error during configuration export", zap.Error(err))
				}
				fmt.Fprintln(os.Stdout, string(btes))
			} else {
				m := flags.AsEnvVariables(conf, upPrefix, true)
				keys := []string{}

				for k := range m {
					keys = append(keys, k)
				}

				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(os.Stdout, "export %s=\"%s\"\n", k, m[k])
				}
			}
		},
	}

	// flags
	configNewCmd.Flags().BoolVar(&configNewAsEnvFlag, "env", false, "Print configuration as environment variable")
	configCmd.AddCommand(configNewCmd)

	// Return base command
	return configCmd
}
