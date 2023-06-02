// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbuilder

import (
	"fmt"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/template/visitor"
	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	"zntr.io/harp/v2/pkg/template/engine"
)

type infrastructureSecretBuilder struct {
	results         chan *bundlev1.Package
	templateContext engine.Context

	// Context
	provider    string
	accountName string
	region      string
	serviceType string
	serviceName string
	err         error
}

// -----------------------------------------------------------------------------

// Infrastructure returns a visitor instance to generate secretpath
// and values.
func infrastructure(results chan *bundlev1.Package, templateContext engine.Context) visitor.InfrastructureVisitor {
	return &infrastructureSecretBuilder{
		results:         results,
		templateContext: templateContext,
	}
}

// -----------------------------------------------------------------------------

func (b *infrastructureSecretBuilder) Error() error {
	return b.err
}

func (b *infrastructureSecretBuilder) VisitForProvider(obj *bundlev1.InfrastructureNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context values
	b.provider, b.err = engine.RenderContext(b.templateContext, obj.Provider)
	if b.err != nil {
		return
	}
	b.accountName, b.err = engine.RenderContext(b.templateContext, obj.Account)
	if b.err != nil {
		return
	}

	// Iterates over regions
	for _, item := range obj.Regions {
		b.VisitForRegion(item)
	}
}

func (b *infrastructureSecretBuilder) VisitForRegion(obj *bundlev1.InfrastructureRegionNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context values
	b.region, b.err = engine.RenderContext(b.templateContext, obj.Name)
	if b.err != nil {
		return
	}

	// Iterates over services
	for _, item := range obj.Services {
		b.VisitForService(item)
	}
}

func (b *infrastructureSecretBuilder) VisitForService(obj *bundlev1.InfrastructureServiceNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context values
	b.serviceType, b.err = engine.RenderContext(b.templateContext, obj.Type)
	if b.err != nil {
		return
	}
	b.serviceName, b.err = engine.RenderContext(b.templateContext, obj.Name)
	if b.err != nil {
		return
	}

	for _, item := range obj.Secrets {
		// Check arguments
		if item == nil {
			continue
		}

		// Parse suffix with template engine
		suffix, err := engine.RenderContext(b.templateContext, item.Suffix)
		if err != nil {
			b.err = fmt.Errorf("unable to merge template is suffix %q", item.Suffix)
			return
		}

		// Generate secret suffix
		secretPath, err := csov1.RingInfra.Path(b.provider, b.accountName, b.region, b.serviceType, b.serviceName, suffix)
		if err != nil {
			b.err = err
			return
		}

		// Prepare template model
		tmplModel := &struct {
			Provider    string
			Account     string
			Region      string
			ServiceType string
			ServiceName string
			Secret      *bundlev1.SecretSuffix
		}{
			Provider:    b.provider,
			Account:     b.accountName,
			Region:      b.region,
			ServiceType: b.serviceType,
			ServiceName: b.serviceName,
			Secret:      item,
		}

		// Compile template
		p, err := parseSecretTemplate(b.templateContext, secretPath, item, tmplModel)
		if err != nil {
			b.err = err
			return
		}

		// Add package to collection
		b.results <- p
	}
}
