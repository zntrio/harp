// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"fmt"
	"strings"
)

// Ring describes secret ring contract.
type Ring interface {
	Level() int
	Name() string
	Prefix() string
	Path(...string) (string, error)
}

var (
	// RingMeta represents R0 secrets.
	RingMeta Ring = &ring{
		level:  0,
		name:   "Meta",
		prefix: ringMeta,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("meta/%s", 1, values...)
		},
	}
	// RingInfra represents R1 secrets.
	RingInfra = &ring{
		level:  1,
		name:   "Infrastructure",
		prefix: ringInfra,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("infra/%s/%s/%s/%s/%s", 5, values...)
		},
	}
	// RingPlatform repsents R2 secrets.
	RingPlatform = &ring{
		level:  2,
		name:   "Platform",
		prefix: ringPlatform,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("platform/%s/%s/%s/%s/%s", 5, values...)
		},
	}
	// RingProduct represents R3 secrets.
	RingProduct = &ring{
		level:  3,
		name:   "Product",
		prefix: ringProduct,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("product/%s/%s/%s/%s", 4, values...)
		},
	}
	// RingApplication represents R4 secrets.
	RingApplication = &ring{
		level:  4,
		name:   "Application",
		prefix: ringApp,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("app/%s/%s/%s/%s/%s/%s", 6, values...)
		},
	}
	// RingArtifact represents R5 secrets.
	RingArtifact = &ring{
		level:  5,
		name:   "Artifact",
		prefix: ringArtifact,
		pathBuilderFunc: func(ring Ring, values ...string) (string, error) {
			return csoPath("artifact/%s/%s/%s", 3, values...)
		},
	}
)

// -----------------------------------------------------------------------------

type ring struct {
	level           int
	name            string
	prefix          string
	pathBuilderFunc func(Ring, ...string) (string, error)
}

func (r ring) Level() int {
	return r.level
}

func (r ring) Name() string {
	return r.name
}

func (r ring) Prefix() string {
	return r.prefix
}

func (r ring) Path(values ...string) (string, error) {
	return r.pathBuilderFunc(r, values...)
}

// -----------------------------------------------------------------------------

// csoPath build and validate a secret path according to CSO specification.
func csoPath(format string, count int, values ...string) (string, error) {
	// Check values count
	if len(values) < count {
		return "", fmt.Errorf("expected (%d) and received (%d) value count doesn't match", count, len(values))
	}

	// Prepare suffix
	suffix := strings.Join(values[count-1:], "/")

	// Prepare values
	var items []interface{}
	for i := 0; i < count-1; i++ {
		items = append(items, values[i])
	}
	items = append(items, suffix)

	// Prepare validation
	csoPath := fmt.Sprintf(format, items...)

	// Validate secret path
	if err := Validate(csoPath); err != nil {
		return "", fmt.Errorf("%q is not a compliant CSO path: %w", csoPath, err)
	}

	// No Error
	return csoPath, nil
}
