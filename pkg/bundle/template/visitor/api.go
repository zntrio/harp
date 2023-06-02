// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package visitor

import (
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// InfrastructureAcceptor is the contract for visitor entrypoint.
type InfrastructureAcceptor interface {
	Accept(InfrastructureVisitor)
}

// PlatformAcceptor is the contract for visitor entrypoint.
type PlatformAcceptor interface {
	Accept(PlatformVisitor)
}

// ProductAcceptor is the contract for visitor entrypoint.
type ProductAcceptor interface {
	Accept(ProductVisitor)
}

// ApplicationAcceptor is the contract for visitor entrypoint.
type ApplicationAcceptor interface {
	Accept(ApplicationVisitor)
}

// InfrastructureVisitor describes visitor method used for tree walk.
type InfrastructureVisitor interface {
	Error() error
	VisitForProvider(obj *bundlev1.InfrastructureNS)
	VisitForRegion(obj *bundlev1.InfrastructureRegionNS)
	VisitForService(obj *bundlev1.InfrastructureServiceNS)
}

// PlatformVisitor describes visitor method used for tree walk.
type PlatformVisitor interface {
	Error() error
	VisitForRegion(obj *bundlev1.PlatformRegionNS)
	VisitForComponent(obj *bundlev1.PlatformComponentNS)
}

// ProductVisitor describes visitor method used for tree walk.
type ProductVisitor interface {
	Error() error
	VisitForComponent(obj *bundlev1.ProductComponentNS)
}

// ApplicationVisitor describes visitor method used for tree walk.
type ApplicationVisitor interface {
	Error() error
	VisitForComponent(obj *bundlev1.ApplicationComponentNS)
}

// TemplateVisitor is a bundle template visitor.
type TemplateVisitor interface {
	Error() error
	Visit(*bundlev1.Template)
}
