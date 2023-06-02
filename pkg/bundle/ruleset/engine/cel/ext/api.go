// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package ext

import (
	"github.com/google/cel-go/cel"
)

var (
	harpPackageObjectType = cel.ObjectType("harp.bundle.v1.Package")
	harpKVObjectType      = cel.ObjectType("harp.bundle.v1.KV")
)
