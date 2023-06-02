// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package reloader

import (
	"context"
	"net"

	"github.com/oklog/run"
)

// Reloader defines socket reloader contract.
type Reloader interface {
	Listen(network, address string) (net.Listener, error)
	SetupGracefulRestart(context.Context, run.Group)
}
