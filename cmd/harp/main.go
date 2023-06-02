// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package main

import (
	"time"

	"zntr.io/harp/v2/cmd/harp/internal/cmd"
	"zntr.io/harp/v2/pkg/sdk/log"
)

func init() {
	// Set default timezone to UTC
	time.Local = time.UTC
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.CheckErr("Unable to complete command execution", err)
	}
}
