// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import (
	"strings"

	csov1 "zntr.io/harp/v2/api/gen/go/cso/v1"
)

const (
	ringMeta     = "meta"
	ringInfra    = "infra"
	ringPlatform = "platform"
	ringProduct  = "product"
	ringApp      = "app"
	ringArtifact = "artifact"
)

// -----------------------------------------------------------------------------

var ringMapNames = strings.Split("invalid;unknown;meta;infra;platform;product;app;artifact", ";")

// ToRingName returns the ring level name.
func ToRingName(lvl csov1.RingLevel) string {
	return ringMapNames[lvl]
}

// FromRingName returns the ring level object according to given name.
func FromRingName(name string) csov1.RingLevel {
	var i int32

	// Search for value
	for idx, n := range ringMapNames {
		if strings.EqualFold(n, name) {
			i = int32(idx)
		}
	}

	return csov1.RingLevel(i)
}

// -----------------------------------------------------------------------------

var qualityMapNames = strings.Split("invalid;unknown;production;staging;qa;dev", ";")

// ToStageName return the stage name.
func ToStageName(lvl csov1.QualityLevel) string {
	return qualityMapNames[lvl]
}

// FromStageName returns the stage level object from given name.
func FromStageName(name string) csov1.QualityLevel {
	var i int32

	// Search for value
	for idx, n := range qualityMapNames {
		if strings.EqualFold(n, name) {
			i = int32(idx)
		}
	}

	return csov1.QualityLevel(i)
}
