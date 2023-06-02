// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package visitor

import (
	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// -----------------------------------------------------------------------------

// InfrastructureDecorator decorates an bundle template infrastructure namespace
// to make it visitable.
func InfrastructureDecorator(entity *bundlev1.InfrastructureNS) InfrastructureAcceptor {
	return &infrastructureDecorator{
		entity: entity,
	}
}

type infrastructureDecorator struct {
	entity *bundlev1.InfrastructureNS
}

func (w *infrastructureDecorator) Accept(v InfrastructureVisitor) {
	v.VisitForProvider(w.entity)
}

// -----------------------------------------------------------------------------

// PlatformDecorator decorates an bundle template platform namespace
// to make it visitable.
func PlatformDecorator(entity *bundlev1.PlatformRegionNS) PlatformAcceptor {
	return &platformDecorator{
		entity: entity,
	}
}

type platformDecorator struct {
	entity *bundlev1.PlatformRegionNS
}

func (w *platformDecorator) Accept(v PlatformVisitor) {
	v.VisitForRegion(w.entity)
}

// -----------------------------------------------------------------------------

// ProductDecorator decorates an bundle template product namespace
// to make it visitable.
func ProductDecorator(entity *bundlev1.ProductComponentNS) ProductAcceptor {
	return &productDecorator{
		entity: entity,
	}
}

type productDecorator struct {
	entity *bundlev1.ProductComponentNS
}

func (w *productDecorator) Accept(v ProductVisitor) {
	v.VisitForComponent(w.entity)
}

// -----------------------------------------------------------------------------

// ApplicationDecorator decorates an bundle template application namespace
// to make it visitable.
func ApplicationDecorator(entity *bundlev1.ApplicationComponentNS) ApplicationAcceptor {
	return &applicationDecorator{
		entity: entity,
	}
}

type applicationDecorator struct {
	entity *bundlev1.ApplicationComponentNS
}

func (w *applicationDecorator) Accept(v ApplicationVisitor) {
	v.VisitForComponent(w.entity)
}
