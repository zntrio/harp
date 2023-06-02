// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"errors"
)

var (
	// ErrPathNotFound is raised when given secret path doesn't exists.
	ErrPathNotFound = errors.New("path not found")
	// ErrSecretNotReadable is raised when trying to read a secret and hit a
	// permission error.
	ErrSecretNotReadable = errors.New("secret not readable")
	// ErrNoData is raised when gievn secret path doesn't contains data.
	ErrNoData = errors.New("no data")
	// ErrCustomMetadataDisabled is raised when trying to write a custom
	// metadata with globally disabled feature.
	ErrCustomMetadataDisabled = errors.New("custom metadata is disabled")
)

// VaultMetadataDataKey represents the secret data key used to store
// metadata.
var VaultMetadataDataKey = "www.vaultproject.io/kv/metadata"

const (
	// CustomMetadataKeyLimit defines the key count limit for custom metadata.
	CustomMetadataKeyLimit = 64
	// CustomMetadataKeySizeLimit defines the key size limit in bytes for
	// custom metadata.
	CustomMetadataKeySizeLimit = 128
	// CustomMetadataValueSizeLimit defines the value size limit in bytes for
	// custom metadata.
	CustomMetadataValueSizeLimit = 512
)

// SecretData is a secret body.
type SecretData map[string]interface{}

// SecretMetadata is secret data attached metadata.
type SecretMetadata map[string]interface{}

// SecretLister repesents secret key listing feature contract.
type SecretLister interface {
	List(ctx context.Context, path string) ([]string, error)
}

// SecretReader represents secret reader feature contract.
type SecretReader interface {
	Read(ctx context.Context, path string) (SecretData, SecretMetadata, error)
	ReadVersion(ctx context.Context, path string, version uint32) (SecretData, SecretMetadata, error)
}

// SecretWriter represents secret writer feature contract.
type SecretWriter interface {
	Write(ctx context.Context, path string, secrets SecretData) error
	WriteWithMeta(ctx context.Context, path string, secrets SecretData, meta SecretMetadata) error
}

// Service declares vault service contract.
type Service interface {
	SecretLister
	SecretReader
	SecretWriter
}
