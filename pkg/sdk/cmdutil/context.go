// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmdutil

import (
	"context"

	"github.com/gosimple/slug"

	"zntr.io/harp/v2/build/version"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// Context initializes a command context.
func Context(ctx context.Context, name string, debug bool, logLevel string) (context.Context, context.CancelFunc) {
	// Context to attach all goroutines
	ctx, cancel := context.WithCancel(ctx)

	// Initialize logger
	log.Setup(ctx,
		&log.Options{
			Debug:    debug,
			LogLevel: logLevel,
			AppName:  slug.Make(name),
			AppID:    version.ID(),
			Version:  version.Version,
			Revision: version.Commit,
		},
	)

	// Return context
	return ctx, cancel
}
