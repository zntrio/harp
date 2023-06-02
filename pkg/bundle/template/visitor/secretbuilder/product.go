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

type productSecretBuilder struct {
	results         chan *bundlev1.Package
	templateContext engine.Context

	// Context
	name      string
	version   string
	component string
	err       error
}

// -----------------------------------------------------------------------------

// Infrastructure returns a visitor instance to generate secretpath
// and values.
func product(results chan *bundlev1.Package, templateContext engine.Context, name, version string) (visitor.ProductVisitor, error) {
	// Parse selector values
	productName, err := engine.RenderContext(templateContext, name)
	if err != nil {
		return nil, fmt.Errorf("unable to render product.name: %w", err)
	}
	if strings.TrimSpace(productName) == "" {
		return nil, fmt.Errorf("product selector must not be empty")
	}

	productVersion, err := engine.RenderContext(templateContext, version)
	if err != nil {
		return nil, fmt.Errorf("unable to render product.version: %w", err)
	}
	if strings.TrimSpace(productVersion) == "" {
		return nil, fmt.Errorf("version selector must not be empty")
	}

	return &productSecretBuilder{
		results:         results,
		templateContext: templateContext,
		name:            productName,
		version:         productVersion,
	}, nil
}

// -----------------------------------------------------------------------------

func (b *productSecretBuilder) Error() error {
	return b.err
}

func (b *productSecretBuilder) VisitForComponent(obj *bundlev1.ProductComponentNS) {
	// Check arguments
	if obj == nil {
		return
	}

	// Set context values
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
		secretPath, err := csov1.RingProduct.Path(b.name, b.version, b.component, suffix)
		if err != nil {
			b.err = err
			return
		}

		// Prepare template model
		tmplModel := &struct {
			Name      string
			Version   string
			Component string
			Secret    *bundlev1.SecretSuffix
		}{
			Name:      b.name,
			Version:   b.version,
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
