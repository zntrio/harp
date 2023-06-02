// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"

	"zntr.io/harp/v2/pkg/vault/logical"
	vpath "zntr.io/harp/v2/pkg/vault/path"
)

type service struct {
	logical   logical.Logical
	mountPath string
	keyName   string
}

// New instantiates a Vault transit backend encryption service.
func New(client *api.Client, mountPath, keyName string) (Service, error) {
	return &service{
		logical:   client.Logical(),
		mountPath: strings.TrimSuffix(path.Clean(mountPath), "/"),
		keyName:   keyName,
	}, nil
}

// -----------------------------------------------------------------------------

func (s *service) Encrypt(ctx context.Context, cleartext []byte) ([]byte, error) {
	// Prepare query
	encryptPath := vpath.SanitizePath(path.Join(url.PathEscape(s.mountPath), "encrypt", url.PathEscape(s.keyName)))
	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(cleartext),
	}

	// Send to Vault.
	secret, err := s.logical.Write(encryptPath, data)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt with %q key: %w", s.keyName, err)
	}

	// Check response wrapping
	if secret.WrapInfo != nil {
		// Unwrap with response token
		secret, err = s.logical.Unwrap(secret.WrapInfo.Token)
		if err != nil {
			return nil, fmt.Errorf("unable to unwrap the response: %w", err)
		}
	}

	// Parse server response.
	if cipherText, ok := secret.Data["ciphertext"].(string); ok && cipherText != "" {
		return []byte(cipherText), nil
	}

	// Return error.
	return nil, errors.New("could not encrypt given data")
}

func (s *service) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	// Prepare query
	decryptPath := vpath.SanitizePath(path.Join(url.PathEscape(s.mountPath), "decrypt", url.PathEscape(s.keyName)))
	data := map[string]interface{}{
		"ciphertext": string(ciphertext),
	}

	// Send to Vault.
	secret, err := s.logical.Write(decryptPath, data)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt with %q key: %w", s.keyName, err)
	}

	// Check response wrapping
	if secret.WrapInfo != nil {
		// Unwrap with response token
		secret, err = s.logical.Unwrap(secret.WrapInfo.Token)
		if err != nil {
			return nil, fmt.Errorf("unable to unwrap the response: %w", err)
		}
	}

	// Parse server response.
	if plainText64, ok := secret.Data["plaintext"].(string); ok && plainText64 != "" {
		plainText, err := base64.StdEncoding.DecodeString(plainText64)
		if err != nil {
			return nil, fmt.Errorf("unable to decode secret: %w", err)
		}

		// Return no error
		return plainText, nil
	}

	// Return error.
	return nil, errors.New("could not decrypt given data")
}
