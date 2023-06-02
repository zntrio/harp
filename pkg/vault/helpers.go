// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// -----------------------------------------------------------------------------

// CheckAuthentication verifies that the connection to vault is setup correctly
// by retrieving information about the configured token.
func CheckAuthentication(ctx context.Context, client *api.Client) ([]string, error) {
	tokenInfo, tokenErr := client.Auth().Token().LookupSelfWithContext(ctx)
	if tokenErr != nil {
		return nil, fmt.Errorf("error connecting to vault: %w", tokenErr)
	}

	tokenPolicies, polErr := tokenInfo.TokenPolicies()
	if polErr != nil {
		return nil, fmt.Errorf("error looking up token policies: %w", tokenErr)
	}

	// No error
	return tokenPolicies, nil
}
