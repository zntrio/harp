// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package platform

import (
	"zntr.io/harp/v2/pkg/sdk/platform/diagnostic"
)

// InstrumentationConfig holds all platform instrumentation settings.
type InstrumentationConfig struct {
	Network    string `toml:"network" default:"tcp" comment:"Network class used for listen (tcp, tcp4, tcp6, unixsocket)"`
	Listen     string `toml:"listen" default:":5556" comment:"Listen address for instrumentation server"`
	Diagnostic struct {
		Enabled bool              `toml:"enabled" default:"false" comment:"Enable diagnostic handlers"`
		Config  diagnostic.Config `toml:"Config" comment:"Diagnostic settings"`
	} `toml:"Diagnostic" comment:"###############################\n Diagnotic Settings \n##############################"`
	Logs struct {
		Level string `toml:"level" default:"warn" comment:"Log level: debug, info, warn, error, dpanic, panic, and fatal"`
	} `toml:"Logs" comment:"###############################\n Logs Settings \n##############################"`
}
