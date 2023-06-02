// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build !fips

package fips

import "os"

func Enabled() bool {
	// Get from env.
	return os.Getenv("HARP_FIPS_MODE") != ""
}
