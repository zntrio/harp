// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"encoding/json"

	"github.com/jmespath/go-jmespath"
	"google.golang.org/protobuf/encoding/protojson"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// MatchJMESPath returns a JMESPatch package matcher specification.
func MatchJMESPath(exp *jmespath.JMESPath) Specification {
	return &jmesPathMatcher{
		exp: exp,
	}
}

type jmesPathMatcher struct {
	exp *jmespath.JMESPath
}

// IsSatisfiedBy returns specification satisfaction status.
func (s *jmesPathMatcher) IsSatisfiedBy(object interface{}) bool {
	// If object is a package
	if p, ok := object.(*bundlev1.Package); ok {
		// Eliminate all package in case of nil query.
		if s.exp == nil {
			return false
		}

		// Rencode as json
		jsonRaw, err := protojson.Marshal(p)
		if err != nil {
			return false
		}

		var object map[string]interface{}
		if errJSON := json.Unmarshal(jsonRaw, &object); errJSON != nil {
			return false
		}

		// Check if query match results
		res, err := s.exp.Search(object)
		if err != nil {
			return false
		}

		// If result is a boolean
		if bRes, ok := res.(bool); ok {
			return bRes
		}
	}

	return false
}
