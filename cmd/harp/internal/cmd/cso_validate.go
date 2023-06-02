// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/log"
)

var (
	csoValidatePaths            []string
	csoValidatePathFrom         string
	csoValidateDropCompliant    bool
	csoValidateDropNonCompliant bool
	csoValidatePathOnly         bool
)

// -----------------------------------------------------------------------------

var csoValidateCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "Validate given paths with CSO Specification",
		Run:     runCSOValidate,
	}

	// Parameters
	cmd.Flags().StringVar(&csoValidatePathFrom, "paths-from", "", "Path to read path from ('-' for stdin or filename)")
	cmd.Flags().StringArrayVar(&csoValidatePaths, "path", []string{}, "Path to validate (multiple)")
	cmd.Flags().BoolVar(&csoValidateDropCompliant, "drop-compliant", false, "Drop compliant path(s) from result")
	cmd.Flags().BoolVar(&csoValidateDropNonCompliant, "drop-non-compliant", false, "Drop non compliant path(s) from result")
	cmd.Flags().BoolVar(&csoValidatePathOnly, "path-only", false, "Display path only as result")

	return cmd
}

type csoValidationResponse struct {
	Compliant bool   `json:"compliant"`
	Error     string `json:"error,omitempty"`
}

func runCSOValidate(cmd *cobra.Command, args []string) {
	ctx, cancel := cmdutil.Context(cmd.Context(), "harp-cso-validate", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	// Check if we have to read external path
	if csoValidatePathFrom != "" {
		// Force read from stdin
		paths, errReader := cmdutil.LineReader(csoValidatePathFrom)
		if errReader != nil {
			log.For(ctx).Fatal("unable to read paths from stdin", zap.Error(errReader))
		}

		// Add to paths
		csoValidatePaths = append(csoValidatePaths, paths...)
	}

	// Check path length
	if len(csoValidatePaths) == 0 {
		log.For(ctx).Fatal("unable to validate empty paths")
	}

	res := map[string]csoValidationResponse{}

	// Validate each path
	for _, p := range csoValidatePaths {
		err := csov1.Validate(p)

		// Error format
		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}

		// Skip result according to parameters
		if csoValidateDropCompliant && err == nil {
			continue
		}
		if csoValidateDropNonCompliant && err != nil {
			continue
		}

		// Add to result
		res[p] = csoValidationResponse{
			Compliant: err == nil,
			Error:     errMessage,
		}
	}

	if !csoValidatePathOnly {
		// Dump as json
		if err := json.NewEncoder(os.Stdout).Encode(res); err != nil {
			log.For(ctx).Fatal("unable to encode validation response", zap.Error(err))
		}
	} else {
		for k := range res {
			fmt.Fprintf(os.Stdout, "%s\n", k)
		}
	}
}
