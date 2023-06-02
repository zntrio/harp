// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package types

import "strings"

// StringArray describes string array type.
type StringArray []string

// -----------------------------------------------------------------------------

// Contains checks if item is in collection.
func (s StringArray) Contains(item string) bool {
	for _, v := range s {
		if strings.EqualFold(item, v) {
			return true
		}
	}

	return false
}

// AddIfNotContains add item if not already in collection.
// Function returns true or false according to add result.
func (s *StringArray) AddIfNotContains(item string) bool {
	if s.Contains(item) {
		// Item not added
		return false
	}
	*s = append(*s, item)

	// Item added
	return true
}

// Remove item from collection.
// Function returns true or false according to removal result.
func (s *StringArray) Remove(item string) bool {
	idx := -1
	for i, v := range *s {
		if strings.EqualFold(item, v) {
			idx = i
			break
		}
	}
	if idx < 0 {
		// Item not removed
		return false
	}
	*s = append((*s)[:idx], (*s)[idx+1:]...)

	// Item removed
	return true
}

// HasOneOf returns true when one of provided items is found in array.
func (s *StringArray) HasOneOf(items ...string) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}

	return false
}
