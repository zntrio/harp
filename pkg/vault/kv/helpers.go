// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"errors"
	"io"

	"github.com/hashicorp/vault/api"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// IsKVv2 detect if the givent path match a kv v2 engine.
func isKVv2(ctx context.Context, secretPath string, client *api.Client) (mountPath string, isV2 bool, err error) {
	mountPath, version, err := kvPreflightVersionRequest(ctx, client, secretPath)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

//nolint:gocyclo,staticcheck // to refactor
func kvPreflightVersionRequest(ctx context.Context, client *api.Client, secretPath string) (mountPath string, backendVersion int, err error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+secretPath)
	resp, err := client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer func(closer io.Closer) {
			log.SafeClose(closer, "unable to successful close request body")
		}(resp.Body)
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, errors.New("nil response from pre-flight request")
	}
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath, ok = mountPathRaw.(string)
		if !ok {
			return "", 0, errors.New("path must be a string")
		}
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	optionMap, ok := options.(map[string]interface{})
	if !ok {
		return mountPath, 1, nil
	}
	versionRaw, hasVersion := optionMap["version"]
	if !hasVersion || versionRaw == nil {
		return mountPath, 1, nil
	}
	version, ok := versionRaw.(string)
	if !ok {
		return "", 0, errors.New("version must be a string")
	}
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}
