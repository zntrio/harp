// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package golang

import (
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// Keep only last 2 versions.
var goVersions = []string{
	"~1.21",
	"~1.20",
}

func init() {
	// Set default timezone to UTC
	time.Local = time.UTC

	if !Is(goVersions...) {
		color.HiRed("#############################################################################################")
		color.HiRed("")
		color.HiRed("Your golang compiler (%s) must be updated to %s to successfully compile all tools.", runtime.Version(), goVersions)
		color.HiRed("")
		color.HiRed("#############################################################################################")
		os.Exit(-1)
	}
}
