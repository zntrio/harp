// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"fmt"
	"strings"

	csov1 "zntr.io/harp/v2/api/gen/go/cso/v1"
)

type ringPacker func([]string) *csov1.Secret

var packMap = map[string]ringPacker{ringMeta: packMeta, ringInfra: packInfra, ringPlatform: packPlatform, ringProduct: packProduct, ringApp: packApplication, ringArtifact: packArtifact}

// Pack a secret path to a protobuf object.
func Pack(secretPath string) (*csov1.Secret, error) {
	// Validate secret path first
	if err := Validate(secretPath); err != nil {
		return nil, fmt.Errorf("unable to pack cso secret: %w", err)
	}

	// Clean path first
	cleanPath := Clean(secretPath)

	// Split path using '/'
	parts := strings.Split(cleanPath, "/")

	// Delegate to ring packer
	rp, ok := packMap[parts[0]]
	if !ok {
		return nil, fmt.Errorf("unable to pack unknown secret %q", parts[0])
	}

	// Call ring packer
	res := rp(parts)

	// No error
	return res, nil
}

// -----------------------------------------------------------------------------

func packMeta(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_META,
		Path: &csov1.Secret_Meta{
			Meta: &csov1.Meta{
				Key: strings.Join(parts[1:], "/"),
			},
		},
	}
}

func packInfra(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_INFRASTRUCTURE,
		Path: &csov1.Secret_Infrastructure{
			Infrastructure: &csov1.Infrastructure{
				CloudProvider: parts[1],
				AccountId:     parts[2],
				Region:        parts[3],
				ServiceName:   parts[4],
				Key:           strings.Join(parts[5:], "/"),
			},
		},
	}
}

func packPlatform(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_PLATFORM,
		Path: &csov1.Secret_Platform{
			Platform: &csov1.Platform{
				Stage:       FromStageName(parts[1]),
				Name:        parts[2],
				Region:      parts[3],
				ServiceName: parts[4],
				Key:         strings.Join(parts[5:], "/"),
			},
		},
	}
}

func packProduct(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_PRODUCT,
		Path: &csov1.Secret_Product{
			Product: &csov1.Product{
				Name:          parts[1],
				Version:       parts[2],
				ComponentName: parts[3],
				Key:           strings.Join(parts[4:], "/"),
			},
		},
	}
}

func packApplication(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_APPLICATION,
		Path: &csov1.Secret_Application{
			Application: &csov1.Application{
				Stage:          FromStageName(parts[1]),
				PlatformName:   parts[2],
				ProductName:    parts[3],
				ProductVersion: parts[4],
				ComponentName:  parts[5],
				Key:            strings.Join(parts[6:], "/"),
			},
		},
	}
}

func packArtifact(parts []string) *csov1.Secret {
	return &csov1.Secret{
		RingLevel: csov1.RingLevel_RING_LEVEL_ARTIFACT,
		Path: &csov1.Secret_Artifact{
			Artifact: &csov1.Artifact{
				Type: parts[1],
				Id:   parts[2],
				Key:  strings.Join(parts[3:], "/"),
			},
		},
	}
}
