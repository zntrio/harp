// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"errors"
	"fmt"

	"github.com/gobwas/glob"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/secret"
)

// KV describes map[string]interface{} alias.
type KV map[string]interface{}

// Glob returns package objects that have name matching the given pattern.
func (kv KV) Glob(pattern string) KV {
	// Prepare Glob filter.
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		g, _ = glob.Compile("**")
	}

	// Apply to collection
	nkv := KV{}
	for name, contents := range kv {
		if g.Match(name) {
			nkv[name] = contents
		}
	}

	return nkv
}

// Get returns a KV of the given package.
func (kv KV) Get(name string) interface{} {
	if v, ok := kv[name]; ok {
		return v
	}
	return KV{}
}

// -----------------------------------------------------------------------------

// AsSecretMap returns a KV map from given package.
func AsSecretMap(p *bundlev1.Package) (KV, error) {
	// Check arguments
	if p == nil {
		return nil, errors.New("unable to transform nil package")
	}

	secrets := KV{}
	for _, s := range p.Secrets.Data {
		// Unpack secret value
		var data interface{}
		if err := secret.Unpack(s.Value, &data); err != nil {
			return nil, fmt.Errorf("unable to unpack %q secret value: %w", p.Name, err)
		}

		// Assign result
		secrets[s.Key] = data
	}

	// No error
	return secrets, nil
}

// FromSecretMap returns the protobuf representation of secretMap.
func FromSecretMap(secretKv KV) ([]*bundlev1.KV, error) {
	secrets := []*bundlev1.KV{}

	// Prepare secret data
	for k, v := range secretKv {
		// Pack secret value
		packed, err := secret.Pack(v)
		if err != nil {
			return nil, fmt.Errorf("unable to pack secret value for `%s`: %w", k, err)
		}

		// Add to secret package
		secrets = append(secrets, &bundlev1.KV{
			Key:   k,
			Type:  fmt.Sprintf("%T", v),
			Value: packed,
		})
	}

	// No error
	return secrets, nil
}
