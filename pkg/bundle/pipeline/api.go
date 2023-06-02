// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package pipeline

import (
	"context"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// Processor declares a bundle processor contract.
type Processor func(context.Context, *bundlev1.Bundle) error

// Context defines tree processing context.
type Context interface {
	GetFile() *bundlev1.Bundle
	GetPackage() *bundlev1.Package
	GetSecret() *bundlev1.SecretChain
	GetKeyValue() *bundlev1.KV
}

// -----------------------------------------------------------------------------

// FileProcessorFunc describes a file object processor contract.
type FileProcessorFunc func(Context, *bundlev1.Bundle) error

// PackageProcessorFunc describes a package object processor contract.
type PackageProcessorFunc func(Context, *bundlev1.Package) error

// ChainProcessorFunc describes a secret chain object processor contract.
type ChainProcessorFunc func(Context, *bundlev1.SecretChain) error

// KVProcessorFunc describes a kv object processor contract.
type KVProcessorFunc func(Context, *bundlev1.KV) error
