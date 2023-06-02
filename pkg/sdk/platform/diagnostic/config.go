// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package diagnostic

// Config holds information for diagnostic handlers.
type Config struct {
	GOPS struct {
		Enabled   bool   `toml:"enabled" default:"false" comment:"Enable GOPS agent"`
		RemoteURL string `toml:"remoteDebugURL" comment:"start a gops agent on specified URL. Ex: localhost:9999"`
	}
	PProf struct {
		Enabled bool `toml:"enabled" default:"true" comment:"Enable PProf handler"`
	}
	ZPages struct {
		Enabled bool `toml:"enabled" default:"true" comment:"Enable zPages handler"`
	}
}

// Validate checks that the configuration is valid.
func (c *Config) Validate() error {
	return nil
}
