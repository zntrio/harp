// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"regexp"

	"github.com/gobwas/glob"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// MatchSecretStrict returns a secret key matcher specification with strict profile.
func MatchSecretStrict(value string) Specification {
	return &matchSecret{
		strict: value,
	}
}

// MatchSecretRegex returns a secret key matcher specification with regexp.
func MatchSecretRegex(regex *regexp.Regexp) Specification {
	return &matchSecret{
		regex: regex,
	}
}

// MatchSecretGlob returns a secret key matcher specification with glob query.
func MatchSecretGlob(pattern string) Specification {
	return &matchPath{
		g: glob.MustCompile(pattern),
	}
}

// matchSecret checks if secret key match the given string.
type matchSecret struct {
	strict string
	regex  *regexp.Regexp
	g      glob.Glob
}

// IsSatisfiedBy returns specification satisfaction status.
func (s *matchSecret) IsSatisfiedBy(object interface{}) bool {
	match := false

	// If object is a package
	if p, ok := object.(*bundlev1.Package); ok {
		// Ignore nil secret package
		if p.Secrets == nil {
			return false
		}

		for _, kv := range p.Secrets.Data {
			switch {
			case s.strict != "":
				return kv.Key == s.strict
			case s.regex != nil:
				return s.regex.MatchString(kv.Key)
			case s.g != nil:
				return s.g.Match(kv.Key)
			}
		}
	}

	return match
}
