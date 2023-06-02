// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package types

import "reflect"

// IsNil returns true if given object is nil.
func IsNil(c interface{}) bool {
	return c == nil ||
		(reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil()) ||
		(reflect.ValueOf(c).Kind() == reflect.Func && reflect.ValueOf(c).IsNil())
}
