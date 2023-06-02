// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"zntr.io/harp/v2/pkg/kv/etcd3"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
	"zntr.io/harp/v2/pkg/sdk/tlsconfig"
	"zntr.io/harp/v2/pkg/tasks/from"
)

// -----------------------------------------------------------------------------

type fromEtcd3Params struct {
	outputPath           string
	basePaths            []string
	lastPathItemAsSecret bool

	endpoints   []string
	dialTimeout time.Duration
	username    string
	password    string

	useTLS             bool
	caFile             string
	certFile           string
	keyFile            string
	passphrase         string
	insecureSkipVerify bool
}

var fromEtcd3Cmd = func() *cobra.Command {
	var params fromEtcd3Params

	cmd := &cobra.Command{
		Use:   "etcd3",
		Short: "Extract KV pairs from CoreOS Etcdv3 KV Store",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger and context
			ctx, cancel := cmdutil.Context(cmd.Context(), "harp-kv-from-etcdv3", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
			defer cancel()

			runFromEtcd3(ctx, &params)
		},
	}

	// Add parameters
	cmd.Flags().StringVar(&params.outputPath, "out", "-", "Container output path ('-' for stdout)")
	cmd.Flags().StringSliceVar(&params.basePaths, "paths", []string{}, "Exported base paths")
	cmd.Flags().BoolVarP(&params.lastPathItemAsSecret, "last-path-item-as-secret-key", "k", false, "Use the last path element as secret key")

	cmd.Flags().StringArrayVar(&params.endpoints, "endpoints", []string{"http://localhost:2379"}, "Etcd cluster endpoints")
	cmd.Flags().DurationVar(&params.dialTimeout, "dial-timeout", 15*time.Second, "Etcd cluster dial timeout")
	cmd.Flags().StringVar(&params.username, "username", "", "Etcd cluster connection username")
	cmd.Flags().StringVar(&params.password, "password", "", "Etcd cluster connection password")

	cmd.Flags().BoolVar(&params.useTLS, "tls", false, "Enable TLS")
	cmd.Flags().StringVar(&params.caFile, "ca-file", "", "TLS CA Certificate file path")
	cmd.Flags().StringVar(&params.certFile, "cert-file", "", "TLS Client certificate file path")
	cmd.Flags().StringVar(&params.keyFile, "key-file", "", "TLS Client private key file path")
	cmd.Flags().StringVar(&params.passphrase, "key-passphrase", "", "TLS Client private key passphrase")
	cmd.Flags().BoolVar(&params.insecureSkipVerify, "insecure-skip-verify", false, "Disable TLS certificate verification")

	return cmd
}

func runFromEtcd3(ctx context.Context, params *fromEtcd3Params) {
	// Create config
	config := clientv3.Config{
		Context:     ctx,
		Endpoints:   params.endpoints,
		DialTimeout: params.dialTimeout,
		Username:    params.username,
		Password:    params.password,
	}

	if params.useTLS {
		tlsConfig, err := tlsconfig.Client(&tlsconfig.Options{
			InsecureSkipVerify: params.insecureSkipVerify,
			CAFile:             params.caFile,
			CertFile:           params.certFile,
			KeyFile:            params.keyFile,
			Passphrase:         params.passphrase,
		})
		if err != nil {
			log.For(ctx).Fatal("unable to initialize TLS settings", zap.Error(err))
			return
		}

		// Assign TLS settings
		config.TLS = tlsConfig
	}

	// Creates a new client
	client, err := clientv3.New(config)
	if err != nil {
		log.For(ctx).Fatal("unable to connect to etcdv3 cluster", zap.Error(err))
		return
	}

	// Prepare store.
	store := etcd3.Store(client)
	defer log.SafeClose(store, "unable to close etcd3 store")

	// Delegate to task
	t := &from.ExtractKVTask{
		Store:                   store,
		ContainerWriter:         cmdutil.FileWriter(params.outputPath),
		BasePaths:               params.basePaths,
		LastPathItemAsSecretKey: params.lastPathItemAsSecret,
	}

	// Run the task
	if err := t.Run(ctx); err != nil {
		log.For(ctx).Fatal("unable to execute kv extract task", zap.Error(err))
		return
	}
}
