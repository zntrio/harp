// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package selector

import (
	"fmt"
	"regexp"

	"github.com/gobwas/glob"
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// MatchPathStrict returns a path matcher specification with strict profile.
func MatchPathStrict(value string) Specification {
	return &matchPath{
		strict: value,
	}
}

// MatchPathRegex returns a path matcher specification with regexp.
func MatchPathRegex(pattern string) (Specification, error) {
	// Compile and check filter
	m, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to compile regex filter: %w", err)
	}

	// No error
	return &matchPath{
		regex: m,
	}, nil
}

// MatchPathGlob returns a path matcher specification with glob query.
func MatchPathGlob(pattern string) (Specification, error) {
	// Compile and check filter
	m, err := glob.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to compile glob filter: %w", err)
	}

	// No error
	return &matchPath{
		g: m,
	}, nil
}

// MatchPath checks if secret path match the given string.
type matchPath struct {
	strict string
	regex  *regexp.Regexp
	g      glob.Glob
}

// IsSatisfiedBy returns specification satisfaction status.
func (s *matchPath) IsSatisfiedBy(object interface{}) bool {
	// If object is a package
	if p, ok := object.(*bundlev1.Package); ok {
		switch {
		case s.strict != "":
			return p.Name == s.strict
		case s.regex != nil:
			return s.regex.MatchString(p.Name)
		case s.g != nil:
			return s.g.Match(p.Name)
		}
	}

	return false
}
