// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package compare

import (
	"fmt"
	"strings"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
)

// ToPatch convert oplog to a bundle patch.
func ToPatch(oplog []DiffItem) (*bundlev1.Patch, error) {
	// Check arguments
	if len(oplog) == 0 {
		return nil, fmt.Errorf("unable to generate a patch with an empty oplog")
	}

	res := &bundlev1.Patch{
		ApiVersion: "harp.zntr.io/v2",
		Kind:       "BundlePatch",
		Meta: &bundlev1.PatchMeta{
			Name:        "autogenerated-patch",
			Description: "Patch generated from oplog",
		},
		Spec: &bundlev1.PatchSpec{
			Rules: []*bundlev1.PatchRule{},
		},
	}

	secretMap := map[string]*bundlev1.PatchRule{}

	// Generate patch rules
	for _, op := range oplog {
		if op.Type == "package" {
			if op.Operation == Remove {
				res.Spec.Rules = append(res.Spec.Rules, &bundlev1.PatchRule{
					Selector: &bundlev1.PatchSelector{
						MatchPath: &bundlev1.PatchSelectorMatchPath{
							Strict: op.Path,
						},
					},
					Package: &bundlev1.PatchPackage{
						Remove: true,
					},
				})
			}
			continue
		}
		if op.Type == "secret" {
			pathParts := strings.SplitN(op.Path, "#", 2)
			pkgRule, ok := secretMap[pathParts[0]]
			if !ok {
				secretMap[pathParts[0]] = &bundlev1.PatchRule{
					Selector: &bundlev1.PatchSelector{
						MatchPath: &bundlev1.PatchSelectorMatchPath{
							Strict: pathParts[0],
						},
					},
					Package: &bundlev1.PatchPackage{
						Data: &bundlev1.PatchSecret{
							Kv: &bundlev1.PatchOperation{},
						},
					},
				}
				pkgRule = secretMap[pathParts[0]]
			}

			switch op.Operation {
			case Add:
				if pkgRule.Package.Data.Kv.Add == nil {
					pkgRule.Package.Data.Kv.Add = map[string]string{}
				}
				pkgRule.Package.Data.Kv.Add[pathParts[1]] = op.Value
			case Replace:
				if pkgRule.Package.Data.Kv.Update == nil {
					pkgRule.Package.Data.Kv.Update = map[string]string{}
				}
				pkgRule.Package.Data.Kv.Update[pathParts[1]] = op.Value
			case Remove:
				if pkgRule.Package.Data.Kv.Remove == nil {
					pkgRule.Package.Data.Kv.Remove = []string{}
				}
				pkgRule.Package.Data.Kv.Remove = append(pkgRule.Package.Data.Kv.Remove, pathParts[1])
			}
		}
	}

	// Add grouped secret patches
	for _, r := range secretMap {
		if r == nil {
			continue
		}

		res.Spec.Rules = append(res.Spec.Rules, r)
	}

	// No error
	return res, nil
}
