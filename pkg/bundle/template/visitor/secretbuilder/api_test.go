// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package secretbuilder

import (
	"testing"

	fuzz "github.com/google/gofuzz"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/template/engine"
)

func TestVisit_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		v := New(&bundlev1.Bundle{}, engine.NewContext())

		tmpl := &bundlev1.Template{
			ApiVersion: "harp.zntr.io/v2",
			Kind:       "BundleTemplate",
			Spec: &bundlev1.TemplateSpec{
				Selector: &bundlev1.Selector{
					Quality:   "production",
					Product:   "harp",
					Version:   "v1.0.0",
					Platform:  "test",
					Component: "cli",
				},
				Namespaces: &bundlev1.Namespaces{},
			},
		}

		// Infrastructure
		f.Fuzz(&tmpl.Spec.Namespaces.Infrastructure)
		// Platform
		f.Fuzz(&tmpl.Spec.Namespaces.Platform)
		// Product
		f.Fuzz(&tmpl.Spec.Namespaces.Product)
		// Application
		f.Fuzz(&tmpl.Spec.Namespaces.Application)

		v.Visit(tmpl)
	}
}

func TestVisit_Template_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		v := New(&bundlev1.Bundle{}, engine.NewContext())

		tmpl := &bundlev1.Template{
			ApiVersion: "harp.zntr.io/v2",
			Kind:       "BundleTemplate",
			Spec: &bundlev1.TemplateSpec{
				Selector: &bundlev1.Selector{
					Quality:   "production",
					Product:   "harp",
					Version:   "v1.0.0",
					Platform:  "test",
					Component: "cli",
				},
				Namespaces: &bundlev1.Namespaces{
					Infrastructure: []*bundlev1.InfrastructureNS{
						{
							Provider: "aws",
							Account:  "test",
							Regions: []*bundlev1.InfrastructureRegionNS{
								{
									Name: "eu-central-1",
									Services: []*bundlev1.InfrastructureServiceNS{
										{
											Name: "ssh",
											Type: "ec2",
											Secrets: []*bundlev1.SecretSuffix{
												{
													Suffix: "test",
												},
											},
										},
									},
								},
							},
						},
					},
					Platform: []*bundlev1.PlatformRegionNS{
						{
							Region: "eu-central-1",
							Components: []*bundlev1.PlatformComponentNS{
								{
									Type: "rds",
									Name: "postgres-1",
									Secrets: []*bundlev1.SecretSuffix{
										{
											Suffix: "test",
										},
									},
								},
							},
						},
					},
					Product: []*bundlev1.ProductComponentNS{
						{
							Type: "service",
							Name: "rest-api",
							Secrets: []*bundlev1.SecretSuffix{
								{
									Suffix: "test",
								},
							},
						},
					},
					Application: []*bundlev1.ApplicationComponentNS{
						{
							Type: "service",
							Name: "web",
							Secrets: []*bundlev1.SecretSuffix{
								{
									Suffix: "test",
								},
							},
						},
					},
				},
			},
		}

		// Infrastructure
		f.Fuzz(&tmpl.Spec.Namespaces.Infrastructure[0].Regions[0].Services[0].Secrets[0].Template)
		// Platform
		f.Fuzz(&tmpl.Spec.Namespaces.Platform[0].Components[0].Secrets[0].Template)
		// Product
		f.Fuzz(&tmpl.Spec.Namespaces.Product[0].Secrets[0].Template)
		// Application
		f.Fuzz(&tmpl.Spec.Namespaces.Application[0].Secrets[0].Template)

		v.Visit(tmpl)
	}
}

func TestVisit_Content_Fuzz(t *testing.T) {
	// Making sure the descrption never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New().NilChance(0).NumElements(1, 5)

		v := New(&bundlev1.Bundle{}, engine.NewContext())

		tmpl := &bundlev1.Template{
			ApiVersion: "harp.zntr.io/v2",
			Kind:       "BundleTemplate",
			Spec: &bundlev1.TemplateSpec{
				Selector: &bundlev1.Selector{
					Quality:   "production",
					Product:   "harp",
					Version:   "v1.0.0",
					Platform:  "test",
					Component: "cli",
				},
				Namespaces: &bundlev1.Namespaces{
					Infrastructure: []*bundlev1.InfrastructureNS{
						{
							Provider: "aws",
							Account:  "test",
							Regions: []*bundlev1.InfrastructureRegionNS{
								{
									Name: "eu-central-1",
									Services: []*bundlev1.InfrastructureServiceNS{
										{
											Name: "ssh",
											Type: "ec2",
											Secrets: []*bundlev1.SecretSuffix{
												{
													Suffix:  "test",
													Content: map[string]string{},
												},
											},
										},
									},
								},
							},
						},
					},
					Platform: []*bundlev1.PlatformRegionNS{
						{
							Region: "eu-central-1",
							Components: []*bundlev1.PlatformComponentNS{
								{
									Type: "rds",
									Name: "postgres-1",
									Secrets: []*bundlev1.SecretSuffix{
										{
											Suffix:  "test",
											Content: map[string]string{},
										},
									},
								},
							},
						},
					},
					Product: []*bundlev1.ProductComponentNS{
						{
							Type: "service",
							Name: "rest-api",
							Secrets: []*bundlev1.SecretSuffix{
								{
									Suffix:  "test",
									Content: map[string]string{},
								},
							},
						},
					},
					Application: []*bundlev1.ApplicationComponentNS{
						{
							Type: "service",
							Name: "web",
							Secrets: []*bundlev1.SecretSuffix{
								{
									Suffix:  "test",
									Content: map[string]string{},
								},
							},
						},
					},
				},
			},
		}

		// Infrastructure
		f.Fuzz(&tmpl.Spec.Namespaces.Infrastructure[0].Regions[0].Services[0].Secrets[0].Content)
		// Platform
		f.Fuzz(&tmpl.Spec.Namespaces.Platform[0].Components[0].Secrets[0].Content)
		// Product
		f.Fuzz(&tmpl.Spec.Namespaces.Product[0].Secrets[0].Content)
		// Application
		f.Fuzz(&tmpl.Spec.Namespaces.Application[0].Secrets[0].Content)

		v.Visit(tmpl)
	}
}
