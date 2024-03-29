// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secret

import (
	"encoding/asn1"
	"fmt"
)

const (
	formatVersion = int(0x00000001)
)

// MustPack uses Pack but panic on error.
func MustPack(value interface{}) []byte {
	out, err := Pack(value)
	if err != nil {
		panic(err)
	}
	return out
}

// Pack a secret value.
func Pack(value interface{}) ([]byte, error) {
	// Encode the payload
	payload, err := asn1.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to pack secret value: %w", err)
	}

	// Pack header
	header, err := asn1.Marshal(formatVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal header of sequence: %w", err)
	}

	// Pack body
	body, err := asn1.Marshal(asn1.RawValue{
		Class:      asn1.ClassUniversal,
		IsCompound: true,
		Tag:        asn1.TagSequence,
		Bytes:      append(header, payload...),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal final sequence: %w", err)
	}

	// No error
	return body, nil
}

// Unpack a secret value.
func Unpack(in []byte, out interface{}) error {
	var raw asn1.RawValue

	_, err := asn1.Unmarshal(in, &raw)
	if err != nil {
		return fmt.Errorf("unable to unpack secret header: %w", err)
	}
	if raw.Class != asn1.ClassUniversal || raw.Tag != asn1.TagSequence || !raw.IsCompound {
		return asn1.StructuralError{Msg: fmt.Sprintf(
			"invalid packed structure object - class [%02x], tag [%02x]",
			raw.Class, raw.Tag)}
	}

	var version int
	rest, err := asn1.Unmarshal(raw.Bytes, &version)
	if err != nil {
		return fmt.Errorf("unable to unpack format version: %w", err)
	}

	// Compare with expected
	if version != formatVersion {
		return fmt.Errorf("unexpected packed version, received %d, expected %d", version, formatVersion)
	}

	// Decode the value
	if _, err := asn1.Unmarshal(rest, out); err != nil {
		return fmt.Errorf("unable to upack secret value: %w", err)
	}

	return nil
}
