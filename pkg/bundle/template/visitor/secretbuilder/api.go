// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbuilder

import (
	"errors"
	"fmt"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/template/visitor"
	"zntr.io/harp/v2/pkg/template/engine"
)

// New returns a secret builder visitor instance.
func New(result *bundlev1.Bundle, templateCtx engine.Context) visitor.TemplateVisitor {
	return &secretBuilder{
		bundle:          result,
		templateContext: templateCtx,
	}
}

// -----------------------------------------------------------------------------

type secretBuilder struct {
	bundle          *bundlev1.Bundle
	templateContext engine.Context
	err             error
}

//nolint:gocyclo,gocognit,funlen // refactoring later
func (sb *secretBuilder) Visit(t *bundlev1.Template) {
	results := make(chan *bundlev1.Package)

	// Check arguments
	if t == nil {
		sb.err = errors.New("template is nil")
		return
	}
	if t.Spec == nil {
		sb.err = errors.New("template spec nil")
		return
	}
	if t.Spec.Namespaces == nil {
		sb.err = errors.New("template spec namespace nil")
		return
	}

	go func() {
		defer close(results)

		// Infrastructure secrets
		if t.Spec.Namespaces.Infrastructure != nil {
			for _, obj := range t.Spec.Namespaces.Infrastructure {
				// Initialize a infrastructure visitor
				v := infrastructure(results, sb.templateContext)

				// Traverse the object-tree
				visitor.InfrastructureDecorator(obj).Accept(v)

				// Get result
				if err := v.Error(); err != nil {
					sb.err = err
					return
				}
			}
		}

		// Platform secrets
		if t.Spec.Namespaces.Platform != nil {
			for _, obj := range t.Spec.Namespaces.Platform {
				// Check selector
				if t.Spec.Selector == nil {
					sb.err = fmt.Errorf("selector is mandatory for platform secrets")
					return
				}

				// Initialize a infrastructure visitor
				v, err := platform(results, sb.templateContext, t.Spec.Selector.Quality, t.Spec.Selector.Platform)
				if err != nil {
					sb.err = err
					return
				}

				// Traverse the object-tree
				visitor.PlatformDecorator(obj).Accept(v)

				// Get result
				if err := v.Error(); err != nil {
					sb.err = err
					return
				}
			}
		}

		// Product secrets
		if t.Spec.Namespaces.Product != nil {
			for _, obj := range t.Spec.Namespaces.Product {
				// Check selector
				if t.Spec.Selector == nil {
					sb.err = fmt.Errorf("selector is mandatory for product secrets")
					return
				}

				// Initialize a infrastructure visitor
				v, err := product(results, sb.templateContext, t.Spec.Selector.Product, t.Spec.Selector.Version)
				if err != nil {
					sb.err = err
					return
				}

				// Traverse the object-tree
				visitor.ProductDecorator(obj).Accept(v)

				// Get result
				if err := v.Error(); err != nil {
					sb.err = err
					return
				}
			}
		}

		// Application secrets
		if t.Spec.Namespaces.Application != nil {
			for _, obj := range t.Spec.Namespaces.Application {
				// Check selector
				if t.Spec.Selector == nil {
					sb.err = fmt.Errorf("selector is mandatory for application secrets")
					return
				}

				// Initialize a infrastructure visitor
				v, err := application(results, sb.templateContext, t.Spec.Selector.Quality, t.Spec.Selector.Platform, t.Spec.Selector.Product, t.Spec.Selector.Version)
				if err != nil {
					sb.err = err
					return
				}

				// Traverse the object-tree
				visitor.ApplicationDecorator(obj).Accept(v)

				// Get result
				if err := v.Error(); err != nil {
					sb.err = err
					return
				}
			}
		}
	}()

	// Pull all packages
	for p := range results {
		sb.bundle.Packages = append(sb.bundle.Packages, p)
	}
}

func (sb *secretBuilder) Error() error {
	return sb.err
}
