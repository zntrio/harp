// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/compare"
	"zntr.io/harp/v2/pkg/bundle/hcl"
	"zntr.io/harp/v2/pkg/bundle/secret"
	"zntr.io/harp/v2/pkg/sdk/types"
)

// FromDump creates a bundle from a JSON Dump.
func FromDump(r io.Reader) (*bundlev1.Bundle, error) {
	// Check parameters
	if types.IsNil(r) {
		return nil, fmt.Errorf("unable to process nil reader")
	}

	// Drain input content
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read input content: %w", err)
	}

	// Build the container from json
	var b bundlev1.Bundle
	if err = protojson.Unmarshal(content, &b); err != nil {
		return nil, fmt.Errorf("unable to decode JSON bundle: %w", err)
	}

	// Convert secret values to current value packing method.
	for _, p := range b.Packages {
		for _, s := range p.Secrets.Data {
			// Decode json encoded value
			var data interface{}
			if errJSON := json.Unmarshal(s.Value, &data); errJSON != nil {
				return nil, fmt.Errorf("unable to decode %q - %q secret value as json: %w", p.Name, s.Key, errJSON)
			}

			// Pack secret value
			payload, err := secret.Pack(data)
			if err != nil {
				return nil, fmt.Errorf("unable to pack %q - %q secret value: %w", p.Name, s.Key, err)
			}

			// Replace current json encoded secret value by packed one.
			s.Value = payload
		}
	}

	// No error
	return &b, nil
}

// FromOpLog convert oplog to a bundle.
func FromOpLog(oplog compare.OpLog) (*bundlev1.Bundle, error) {
	// Create an empty bundle.
	b := &bundlev1.Bundle{}

	packageMap := map[string]*bundlev1.Package{}

	// Generate patch rules
	for _, op := range oplog {
		switch op.Type {
		case "package":
			// Ignore package operation
			continue
		case "secret":
			pathParts := strings.SplitN(op.Path, "#", 2)
			pkg, ok := packageMap[pathParts[0]]
			if !ok {
				packageMap[pathParts[0]] = &bundlev1.Package{
					Name: pathParts[0],
					Secrets: &bundlev1.SecretChain{
						Data: []*bundlev1.KV{},
					},
				}
				pkg = packageMap[pathParts[0]]
			}

			// Process oplog event
			switch op.Operation {
			case compare.Add, compare.Replace:
				// Pack secret value
				payload, err := secret.Pack(op.Value)
				if err != nil {
					return nil, fmt.Errorf("unable to pack secret value for %q / %q: %w", pathParts[0], pathParts[1], err)
				}

				// Assign secret data
				pkg.Secrets.Data = append(pkg.Secrets.Data, &bundlev1.KV{
					Key:   pathParts[1],
					Type:  "string",
					Value: payload,
				})
			case compare.Remove:
				// Ignore secret removal
			}
		default:
			return nil, fmt.Errorf("unknown oplog type %q", op.Type)
		}
	}

	// Assign packages
	for _, p := range packageMap {
		b.Packages = append(b.Packages, p)
	}

	// No error
	return b, nil
}

// FromMap builds a secret container from map K/V.
func FromMap(input map[string]KV) (*bundlev1.Bundle, error) {
	// Check input
	if input == nil {
		return nil, fmt.Errorf("unable to process nil map")
	}

	res := &bundlev1.Bundle{
		Packages: []*bundlev1.Package{},
	}
	for packageName, secretKv := range input {
		// Prepare a package
		p := &bundlev1.Package{
			Name:    packageName,
			Secrets: &bundlev1.SecretChain{},
		}

		// Prepare secret data
		for k, v := range secretKv {
			// Pack secret value
			packed, err := secret.Pack(v)
			if err != nil {
				return nil, fmt.Errorf("unable to pack secret value for `%s`: %w", fmt.Sprintf("%s.%s", packageName, k), err)
			}

			// Add to secret package
			p.Secrets.Data = append(p.Secrets.Data, &bundlev1.KV{
				Key:   k,
				Type:  fmt.Sprintf("%T", v),
				Value: packed,
			})
		}

		// Add package to result
		res.Packages = append(res.Packages, p)
	}

	// No error
	return res, nil
}

// FromHCL convert HCL-DSL to a bundle.
func FromHCL(input *hcl.Config) (*bundlev1.Bundle, error) {
	// Check arguments
	if input == nil {
		return nil, errors.New("unable to process nil hcl object")
	}

	// Create an empty bundle.
	res := &bundlev1.Bundle{
		Labels:      input.Labels,
		Annotations: input.Annotations,
		Packages:    []*bundlev1.Package{},
	}

	for _, pkg := range input.Packages {
		// Prepare a package
		p := &bundlev1.Package{
			Name:        pkg.Path,
			Annotations: pkg.Annotations,
			Labels:      pkg.Labels,
			Secrets:     &bundlev1.SecretChain{},
		}

		// Prepare secret data
		for k, v := range pkg.Secrets {
			// Pack secret value
			packed, err := secret.Pack(v)
			if err != nil {
				return nil, fmt.Errorf("unable to pack secret value for `%s`: %w", fmt.Sprintf("%s.%s", pkg.Path, k), err)
			}

			// Add to secret package
			p.Secrets.Data = append(p.Secrets.Data, &bundlev1.KV{
				Key:   k,
				Type:  fmt.Sprintf("%T", v),
				Value: packed,
			})
		}

		// Add package to result
		res.Packages = append(res.Packages, p)
	}

	// No error
	return res, nil
}
