// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"strings"
)

func Clean(secretPath string) string {
	// Remove / in prefix
	s := strings.TrimPrefix(secretPath, "/")

	// Remove spaces
	s = strings.TrimSpace(s)

	// Lowercase everything
	s = strings.ToLower(s)

	// Return secret path
	return s
}
