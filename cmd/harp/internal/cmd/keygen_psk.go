// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/awnumar/memguard"
	"github.com/spf13/cobra"

	"github.com/zntrio/harp/v2/pkg/sdk/cmdutil"
)

// -----------------------------------------------------------------------------

var keygenPreSharedKeyCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pre-shared-key",
		Aliases: []string{"psk"},
		Short:   "Generate and print a container pre-shared-key",
		Run:     runKeygenPreSharedKey,
	}

	return cmd
}

func runKeygenPreSharedKey(cmd *cobra.Command, args []string) {
	_, cancel := cmdutil.Context(cmd.Context(), "harp-keygen-psk", conf.Debug.Enabled, conf.Instrumentation.Logs.Level)
	defer cancel()

	fmt.Fprintf(os.Stdout, "%s", base64.RawURLEncoding.EncodeToString(memguard.NewBufferRandom(64).Bytes()))
}
