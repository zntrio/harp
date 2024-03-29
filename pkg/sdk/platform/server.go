// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package platform

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dchest/uniuri"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/platform/diagnostic"
	"zntr.io/harp/v2/pkg/sdk/platform/reloader"
)

// -----------------------------------------------------------------------------

// Server represents platform server.
type Server struct {
	Debug           bool
	Name            string
	Version         string
	Revision        string
	Instrumentation InstrumentationConfig
	Network         string
	Address         string
	Builder         func(ln net.Listener, group *run.Group)
}

// Serve starts the server listening process.
func Serve(ctx context.Context, srv *Server) error {
	// Generate an instance identifier
	appID := uniuri.NewLen(64)

	// Prepare logger
	log.Setup(ctx, &log.Options{
		Debug:    srv.Debug,
		AppName:  srv.Name,
		AppID:    appID,
		Version:  srv.Version,
		Revision: srv.Revision,
		LogLevel: srv.Instrumentation.Logs.Level,
	})

	// Preparing instrumentation
	instrumentationRouter := instrumentServer(ctx, srv)

	// Configure graceful restart
	upg := reloader.Create(ctx)

	var group run.Group

	// Instrumentation server
	{
		ln, err := upg.Listen(srv.Instrumentation.Network, srv.Instrumentation.Listen)
		if err != nil {
			return fmt.Errorf("platform: unable to start instrumentation server: %w", err)
		}

		server := &http.Server{
			Handler: instrumentationRouter,
			// Set timeouts to avoid Slowloris attacks.
			ReadHeaderTimeout: time.Second * 20,
			WriteTimeout:      time.Second * 60,
			ReadTimeout:       time.Second * 60,
			IdleTimeout:       time.Second * 120,
		}

		group.Add(
			func() error {
				log.For(ctx).Info("Starting instrumentation server", zap.String("address", ln.Addr().String()))
				return server.Serve(ln)
			},
			func(e error) {
				log.For(ctx).Info("Shutting instrumentation server down")

				ctxShutdown, cancel := context.WithTimeout(ctx, 60*time.Second)
				defer cancel()

				log.CheckErrCtx(ctx, "Error raised while shutting down the server", server.Shutdown(ctxShutdown))
				log.SafeClose(server, "Unable to close instrumentation server")
			},
		)
	}

	// Initialiaze network listener
	ln, err := upg.Listen(srv.Network, srv.Address)
	if err != nil {
		return fmt.Errorf("unable to start server listener: %w", err)
	}

	// Initialize the component
	srv.Builder(ln, &group)

	// Setup signal handler
	{
		var (
			cancelInterrupt = make(chan struct{})
			ch              = make(chan os.Signal, 2)
		)
		defer close(ch)

		group.Add(
			func() error {
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

				select {
				case <-ch:
					log.For(ctx).Info("Captured signal")
				case <-cancelInterrupt:
				}

				return nil
			},
			func(e error) {
				close(cancelInterrupt)
				signal.Stop(ch)
			},
		)
	}

	// Register graceful restart handler
	upg.SetupGracefulRestart(ctx, group)

	// Run goroutine group
	return group.Run()
}

func instrumentServer(ctx context.Context, srv *Server) *http.ServeMux {
	instrumentationRouter := http.NewServeMux()

	// Register common features
	if srv.Instrumentation.Diagnostic.Enabled {
		cancelFunc, err := diagnostic.Register(ctx, &srv.Instrumentation.Diagnostic.Config, instrumentationRouter)
		if err != nil {
			log.For(ctx).Fatal("Unable to register diagnostic instrumentation", zap.Error(err))
		}
		defer cancelFunc()
	}

	return instrumentationRouter
}
