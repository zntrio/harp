// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbuilder

import (
	"fmt"
	"strings"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/template/visitor"
	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	"zntr.io/harp/v2/pkg/template/engine"
)

type platformSecretBuilder struct {
	results         chan *bundlev1.Package
	templateContext engine.Context

	// Context
	quality   string
	name      string
	region    string
	component string
	err       error
}

// -----------------------------------------------------------------------------

// Infrastructure returns a visitor instance to generate secretpath
// and values.
func platform(results chan *bundlev1.Package, templateContext engine.Context, quality, name string) (visitor.PlatformVisitor, error) {
	// Parse selector values
	platformQuality, err := engine.RenderContext(templateContext, quality)
	if err != nil {
		return nil, fmt.Errorf("unable to render platform.quality: %w", err)
	}
	if strings.TrimSpace(platformQuality) == "" {
		return nil, fmt.Errorf("quality selector must not be empty")
	}

	platformName, err := engine.RenderContext(templateContext, name)
	if err != nil {
		return nil, fmt.Errorf("unable to render platform.name: %w", err)
	}
	if strings.TrimSpace(platformName) == "" {
		return nil, fmt.Errorf("platform selector must not be empty")
	}

	return &platformSecretBuilder{
		results:         results,
		templateContext: templateContext,
		quality:         platformQuality,
		name:            platformName,
	}, nil
}

// -----------------------------------------------------------------------------

func (b *platformSecretBuilder) Error() error {
	return b.err
}

func (b *platformSecretBuilder) VisitForRegion(obj *bundlev1.PlatformRegionNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context values
	b.region, b.err = engine.RenderContext(b.templateContext, obj.Region)
	if b.err != nil {
		return
	}

	// Iterates over all components
	for _, item := range obj.Components {
		b.VisitForComponent(item)
	}
}

func (b *platformSecretBuilder) VisitForComponent(obj *bundlev1.PlatformComponentNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context value
	b.component, b.err = engine.RenderContext(b.templateContext, obj.Name)
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
		secretPath, err := csov1.RingPlatform.Path(b.quality, b.name, b.region, b.component, suffix)
		if err != nil {
			b.err = err
			return
		}

		// Prepare template model
		tmplModel := &struct {
			Quality   string
			Name      string
			Region    string
			Component string
			Secret    *bundlev1.SecretSuffix
		}{
			Quality:   b.quality,
			Name:      b.name,
			Region:    b.region,
			Component: b.component,
			Secret:    item,
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
