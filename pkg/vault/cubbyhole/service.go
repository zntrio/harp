// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cubbyhole

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/golang/snappy"
	"github.com/hashicorp/vault/api"

	"zntr.io/harp/v2/pkg/vault/logical"
	vpath "zntr.io/harp/v2/pkg/vault/path"
)

type service struct {
	logical   logical.Logical
	mountPath string
}

// New instantiates a Vault cubbyhole backend service.
func New(client *api.Client, mountPath string) (Service, error) {
	// Apply default cubbyhole mountpath if not overrided.
	if mountPath == "" {
		mountPath = "cubbyhole"
	}

	return &service{
		logical:   client.Logical(),
		mountPath: vpath.SanitizePath(mountPath),
	}, nil
}

// -----------------------------------------------------------------------------

const (
	secretSizeLimit = 1024 * 1024 * 1024 // 1Mb
)

// Put a secret in cubbyhole to retrieve a wrapping token.
//
//nolint:interfacer // -- wants to replace time.Duration by fmt.Stringer
func (s *service) Put(_ context.Context, r io.Reader) (string, error) {
	// Encode secret
	payload, err := io.ReadAll(io.LimitReader(r, secretSizeLimit))
	if err != nil {
		return "", fmt.Errorf("unable to drain secret reader: %w", err)
	}

	// Compress and encode
	final := base64.StdEncoding.EncodeToString(snappy.Encode(nil, payload))

	// Add to cubbyhole
	return addToCubbyhole(s.logical, s.mountPath, final)
}

// Get a secret from wrapping token.
func (s *service) Get(_ context.Context, token string, w io.Writer) error {
	// Unwrap token
	encoded, err := unWrap(s.logical, token)
	if err != nil {
		return err
	}

	// Decode
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("invalid secret payload: %w", err)
	}

	// Decompress
	final, err := snappy.Decode(nil, payload)
	if err != nil {
		return fmt.Errorf("invalid secret payload: %w", err)
	}

	// Return result to writer
	_, err = w.Write(final)
	if err != nil {
		return fmt.Errorf("unable to write result to the writer: %w", err)
	}

	// No error
	return nil
}
