// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package main

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/pipeline"
	"zntr.io/harp/v2/pkg/sdk/log"
)

func main() {
	var (
		// Initialize an execution context
		ctx = context.Background()
	)

	// Run the pipeline
	if err := pipeline.Run(ctx,
		pipeline.PackageProcessor(packageRemapper), // Package processor
	); err != nil {
		log.For(ctx).Fatal("unable to process bundle", zap.Error(err))
	}
}

// -----------------------------------------------------------------------------

func packageRemapper(ctx pipeline.Context, p *bundlev1.Package) error {

	// Remapping condition
	if !strings.HasPrefix(p.Name, "services/production/global/clusters/") {
		// Skip path remapping
		return nil
	}

	// Remap secret path
	p.Name = fmt.Sprintf("app/production/global/clusters/v1.0.0/bootstrap/%s", strings.TrimPrefix(p.Name, "services/production/global/clusters/"))

	// No error
	return nil
}
