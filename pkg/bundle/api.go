// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"io"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// Reader exposes bundle reader contract.
type Reader interface {
	Read(reader io.Reader) (*bundlev1.Bundle, error)
}

// Writer exposes bundle writer contract.
type Writer interface {
	Write(file *bundlev1.Bundle) error
}

// Visitor delares the bundle vistor contract.
type Visitor interface {
	Error() error
	VisitForFile(obj *bundlev1.Bundle)
	VisitForPackage(obj *bundlev1.Package)
	VisitForChain(obj *bundlev1.SecretChain)
	VisitForKV(obj *bundlev1.KV)
}

// -----------------------------------------------------------------------------

var (
	bundleAnnotationsKey        = "harp.zntr.io/v2/bundle#annotations"
	bundleLabelsKey             = "harp.zntr.io/v2/bundle#labels"
	packageAnnotations          = "harp.zntr.io/v2/package#annotations"
	packageLabels               = "harp.zntr.io/v2/package#labels"
	packageEncryptionAnnotation = "harp.zntr.io/v2/package#encryptionKeyAlias"
	packageEncryptedValueType   = "harp.zntr.io/v2/package#encryptedValue"
)

// AnnotationOwner defines annotations owner contract.
type AnnotationOwner interface {
	GetAnnotations() map[string]string
}

// LabelOwner defines label owner contract.
type LabelOwner interface {
	GetLabels() map[string]string
}
