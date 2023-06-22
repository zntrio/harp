// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package config

import (
	"fmt"
	"os"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/flags"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// Load a config
// Apply defaults first, then environment, then finally config file.
func Load(conf interface{}, envPrefix, cfgFile string) error {
	// Apply defaults first
	defaults.SetDefaults(conf)

	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// Overrides with environment
	for k := range flags.AsEnvVariables(conf, "", false) {
		envName := fmt.Sprintf("%s_%s", upPrefix, k)
		log.CheckErr("unable to bind environment variable", viper.BindEnv(strings.ToLower(strings.ReplaceAll(k, "_", ".")), envName), zap.String("var", envName))
	}

	// Apply file settings
	if cfgFile != "" {
		// If the config file doesn't exists, let's exit
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return fmt.Errorf("config: unable to open non-existing file %q: %w", cfgFile, err)
		}

		log.Bg().Info("Load settings from file", zap.String("path", cfgFile))

		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("config: unable to decode config file %q: %w", cfgFile, err)
		}
	}

	// Update viper values
	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("config: unable to apply config %q: %w", cfgFile, err)
	}

	// No error
	return nil
}
