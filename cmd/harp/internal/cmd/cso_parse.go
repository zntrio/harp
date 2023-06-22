// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
)

var (
	csoParsePath   string
	csoParseAsText bool
)

// -----------------------------------------------------------------------------

var csoParseCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse",
		Aliases: []string{"p"},
		Short:   "Parse given CSO path",
		Run:     runCSOParse,
	}

	// Parameters
	cmd.Flags().StringVar(&csoParsePath, "path", "", "Path to parse")
	cmd.Flags().BoolVar(&csoParseAsText, "text", false, "Display path component as text")

	return cmd
}

func runCSOParse(cmd *cobra.Command, args []string) {
	ctx, cancel := cmdutil.Context(cmd.Context(), "harp-cso-parse", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// Validate and pack secret path first
	s, err := csov1.Pack(csoParsePath)
	if err != nil {
		log.For(ctx).Fatal("unable to validate given path as a compliant CSO path", zap.Error(err), zap.String("path", csoParsePath))
	}

	if csoParseAsText {
		if err := csov1.Interpret(s, csov1.Text(), os.Stdout); err != nil {
			log.For(ctx).Fatal("unable to generate textual interpretation of given path", zap.Error(err), zap.String("path", csoParsePath))
		}
	} else {
		// Override values as nil
		s.Value = nil

		// Marshal using protojson
		out, err := protojson.Marshal(s)
		if err != nil {
			log.For(ctx).Fatal("unable to generate json interpretation of given path", zap.Error(err), zap.String("path", csoParsePath))
		}

		// Dump in stdout
		fmt.Fprintf(os.Stdout, "%s", string(out))
	}
}
