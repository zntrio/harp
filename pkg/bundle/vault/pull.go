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
	"zntr.io/harp/v2/pkg/vault/kv"
	vpath "zntr.io/harp/v2/pkg/vault/path"

	"golang.org/x/sync/errgroup"
)

// Pull all given path as a bundle.
func Pull(ctx context.Context, client *api.Client, paths []string, opts ...Option) (*bundlev1.Bundle, error) {
	// Check parameters
	if client == nil {
		return nil, fmt.Errorf("unable to process with nil client")
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no path given to pull")
	}

	// Default values
	var (
		defaultPrefix             = ""
		defaultPathInclusions     = []*regexp.Regexp{}
		defaultPathExclusions     = []*regexp.Regexp{}
		defaultWithSecretMetadata = false
		defaultWithVaultMetadata  = false
		defaultWorkerCount        = int64(4)
		defaultContinueOnError    = false
	)

	// Create default option instance
	defaultOpts := &options{
		prefix:             defaultPrefix,
		exclusions:         defaultPathExclusions,
		includes:           defaultPathInclusions,
		withSecretMetadata: defaultWithSecretMetadata,
		withVaultMetadata:  defaultWithVaultMetadata,
		workerCount:        defaultWorkerCount,
		continueOnError:    defaultContinueOnError,
	}

	// Apply option functions
	for _, o := range opts {
		if err := o(defaultOpts); err != nil {
			return nil, fmt.Errorf("unable to apply option %T: %w", o, err)
		}
	}

	// Run the pull process
	b, err := runPull(ctx, client, paths, defaultOpts)
	if err != nil {
		return nil, fmt.Errorf("error occurs during pull process: %w", err)
	}

	// No error
	return b, nil
}

// runPull starts a multithreaded Vault secret puller.
func runPull(ctx context.Context, client *api.Client, paths []string, opts *options) (*bundlev1.Bundle, error) {
	var res *bundlev1.Bundle

	// Initialize operation
	packageChan := make(chan *bundlev1.Package)

	// Prepare output
	g, gctx := errgroup.WithContext(ctx)

	// Preprocess paths
	if len(opts.exclusions) > 0 {
		paths = collect(paths, opts.exclusions, false)
	}
	if len(opts.includes) > 0 {
		paths = collect(paths, opts.includes, true)
	}

	// Fork consumer

	// Secret packages consumer
	g.Go(func() error {
		b := &bundlev1.Bundle{}

		// Wait for all packages
		for p := range packageChan {
			b.Packages = append(b.Packages, p)
		}

		// Assign result
		res = b

		// No error
		return nil
	})

	// Fork reader
	g.Go(func() error {
		defer close(packageChan)

		gReader, gReaderctx := errgroup.WithContext(gctx)

		// Wrap process in a builder to be able to pass p parameter
		exportBuilder := func(p string) func() error {
			return func() error {
				// Create dedicated service reader
				service, err := kv.New(client, p, kv.WithVaultMetatadata(opts.withVaultMetadata), kv.WithContext(gReaderctx))
				if err != nil {
					return fmt.Errorf("unable to prepare vault reader for path %q: %w", p, err)
				}

				// Create an exporter
				op := operation.Exporter(service, vpath.SanitizePath(p), packageChan, opts.withSecretMetadata, opts.workerCount, opts.continueOnError)

				// Run the job
				if err := op.Run(gReaderctx); err != nil {
					return fmt.Errorf("unable to export secret values for path `%s': %w", p, err)
				}

				// No error
				return nil
			}
		}

		// Generate producers
		for _, p := range paths {
			// For the process
			gReader.Go(exportBuilder(p))
		}

		// Wait for all producers to finish
		if err := gReader.Wait(); err != nil {
			return fmt.Errorf("unable to read secrets: %w", err)
		}

		// No error
		return nil
	})

	// Wait for completion
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("unable to pull secrets: %w", err)
	}

	// Check bundle result
	if res == nil {
		return nil, fmt.Errorf("result bundle is nil")
	}

	// No error
	return res, nil
}
