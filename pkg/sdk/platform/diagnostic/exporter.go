// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package diagnostic

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/google/gops/agent"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// Register adds diagnostic tools to main process.
func Register(ctx context.Context, conf *Config, r *http.ServeMux) (func(), error) {
	if conf.GOPS.Enabled {
		// Start diagnostic handler
		if conf.GOPS.RemoteURL != "" {
			log.For(ctx).Info("Starting gops agent", zap.String("url", conf.GOPS.RemoteURL))

			if err := agent.Listen(agent.Options{Addr: conf.GOPS.RemoteURL}); err != nil {
				log.For(ctx).Error("Error on starting gops agent", zap.Error(err))
			}
		} else {
			log.For(ctx).Info("Starting gops agent locally")

			if err := agent.Listen(agent.Options{}); err != nil {
				log.For(ctx).Error("Error on starting gops agent locally", zap.Error(err))
			}
		}
	}

	if conf.PProf.Enabled {
		r.HandleFunc("/debug/pprof", pprof.Index)
		r.HandleFunc("/debug/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/profile", pprof.Profile)
		r.HandleFunc("/debug/symbol", pprof.Symbol)
		r.HandleFunc("/debug/trace", pprof.Trace)
		r.Handle("/debug/goroutine", pprof.Handler("goroutine"))
		r.Handle("/debug/heap", pprof.Handler("heap"))
		r.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/debug/block", pprof.Handler("block"))
	}

	// No error
	return func() {
		agent.Close()
	}, nil
}
