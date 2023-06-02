// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

//go:build windows
// +build windows

package reloader

import (
	"context"
	"net"

	"github.com/oklog/run"

	"zntr.io/harp/v2/pkg/sdk/log"
)

// UnsupportedReloader is the file descriptor reloader mock for Windows.
type UnsupportedReloader struct {
}

// Create a descriptor reloader.
func Create(ctx context.Context) Reloader {
	log.For(ctx).Warn("graceful reload is not supported on this platform")
	return &UnsupportedReloader{}
}

// Listen create a listener socket.
func (t *UnsupportedReloader) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

// SetupGracefulRestart does nothing on Windows.
func (t *UnsupportedReloader) SetupGracefulRestart(context context.Context, group run.Group) {
	// no-op since it isn't supported
}
