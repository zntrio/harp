// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/vault/api"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/vault/internal/operation"
)

// Push the given bundle in Hashicorp Vault.
func Push(ctx context.Context, b *bundlev1.Bundle, client *api.Client, opts ...Option) error {
	// Check parameters
	if b == nil {
		return fmt.Errorf("unable to process nil bundle")
	}
	if client == nil {
		return fmt.Errorf("unable to process nil vault client")
	}

	// Default values
	var (
		defaultPrefix             = ""
		defaultPathInclusions     = []*regexp.Regexp{}
		defaultPathExclusions     = []*regexp.Regexp{}
		defaultWithSecretMetadata = false
		defaultWithVaultMetadata  = false
		defaultWorkerCount        = int64(4)
	)

	// Create default option instance
	defaultOpts := &options{
		prefix:             defaultPrefix,
		exclusions:         defaultPathExclusions,
		includes:           defaultPathInclusions,
		withSecretMetadata: defaultWithSecretMetadata,
		withVaultMetadata:  defaultWithVaultMetadata,
		workerCount:        defaultWorkerCount,
	}

	// Apply option functions
	for _, o := range opts {
		if err := o(defaultOpts); err != nil {
			return fmt.Errorf("unable to apply option: %w", err)
		}
	}

	// No error
	return runPush(ctx, b, client, defaultOpts)
}

func runPush(ctx context.Context, b *bundlev1.Bundle, client *api.Client, opts *options) error {
	// Prepare bundle
	if len(opts.includes) > 0 {
		filteredPackages := []*bundlev1.Package{}
		for _, p := range b.Packages {
			if matchPathRule(p.Name, opts.exclusions) {
				filteredPackages = append(filteredPackages, p)
			}
		}
		b.Packages = filteredPackages
	}
	if len(opts.exclusions) > 0 {
		filteredPackages := []*bundlev1.Package{}
		for _, p := range b.Packages {
			if !matchPathRule(p.Name, opts.exclusions) {
				filteredPackages = append(filteredPackages, p)
			}
		}
		b.Packages = filteredPackages
	}

	// Initialize operation
	op := operation.Importer(client, b, opts.prefix, opts.withSecretMetadata, opts.withVaultMetadata, opts.workerCount)

	// Run the vault operation
	if err := op.Run(ctx); err != nil {
		return fmt.Errorf("unable to push secret bundle: %w", err)
	}

	// No error
	return nil
}
