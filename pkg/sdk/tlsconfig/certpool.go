// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package tlsconfig

import (
	"crypto/x509"
)

// SystemCertPool returns an new empty cert pool,
// accessing system cert pool.
func SystemCertPool() (*x509.CertPool, error) {
	return x509.NewCertPool(), nil
}
