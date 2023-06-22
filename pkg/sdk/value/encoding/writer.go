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
	"zntr.io/harp/v2/pkg/sdk/ioutil"
)

// -----------------------------------------------------------------------------

// NewWriter returns the appropriate writer implementation according to given encoding.
func NewWriter(w io.Writer, encoding string) (io.WriteCloser, error) {
	// Normalize input
	encoding = strings.TrimSpace(strings.ToLower(encoding))

	var encoderWriter io.WriteCloser

	// Apply transformation
	switch encoding {
	case "identity":
		encoderWriter = ioutil.NopCloserWriter(w)
	case "hex", "base16":
		encoderWriter = ioutil.NopCloserWriter(hex.NewEncoder(w))
	case "base32":
		encoderWriter = base32.NewEncoder(base32.StdEncoding, w)
	case "base32hex":
		encoderWriter = base32.NewEncoder(base32.HexEncoding, w)
	case "base62":
		encoderWriter = base62.NewEncoder(base62.StdEncoding, w)
	case "base64":
		encoderWriter = base64.NewEncoder(base64.StdEncoding, w)
	case "base64raw":
		encoderWriter = base64.NewEncoder(base64.RawStdEncoding, w)
	case "base64url":
		encoderWriter = base64.NewEncoder(base64.URLEncoding, w)
	case "base64urlraw":
		encoderWriter = base64.NewEncoder(base64.RawURLEncoding, w)
	case "base85":
		encoderWriter = ascii85.NewEncoder(w)
	default:
		return nil, fmt.Errorf("unhandled encoding strategy %q", encoding)
	}

	// No error
	return encoderWriter, nil
}
