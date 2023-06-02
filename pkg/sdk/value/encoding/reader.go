// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package encoding

import (
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/lytics/base62"
)

// -----------------------------------------------------------------------------

// NewReader returns a reader implementation matching the given encoding strategy.
func NewReader(r io.Reader, encoding string) (io.Reader, error) {
	// Normalize input
	encoding = strings.TrimSpace(strings.ToLower(encoding))

	var decoderReader io.Reader

	// Apply transformation
	switch encoding {
	case "identity":
		decoderReader = r
	case "hex", "base16":
		decoderReader = hex.NewDecoder(r)
	case "base32":
		decoderReader = base32.NewDecoder(base32.StdEncoding, r)
	case "base32hex":
		decoderReader = base32.NewDecoder(base32.HexEncoding, r)
	case "base62":
		decoderReader = base62.NewDecoder(base62.StdEncoding, r)
	case "base64":
		decoderReader = base64.NewDecoder(base64.StdEncoding, r)
	case "base64raw":
		decoderReader = base64.NewDecoder(base64.RawStdEncoding, r)
	case "base64url":
		decoderReader = base64.NewDecoder(base64.URLEncoding, r)
	case "base64urlraw":
		decoderReader = base64.NewDecoder(base64.RawURLEncoding, r)
	case "base85":
		decoderReader = ascii85.NewDecoder(r)
	default:
		return nil, fmt.Errorf("unhandled decoding strategy %q", encoding)
	}

	// No error
	return decoderReader, nil
}
