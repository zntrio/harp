// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ext

import (
	"reflect"
	"strings"

	"github.com/gobwas/glob"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	csov1 "zntr.io/harp/v2/pkg/cso/v1"
	htypes "zntr.io/harp/v2/pkg/sdk/types"
)

// Packages exported package operations.
func Packages() cel.EnvOption {
	return cel.Lib(packageLib{})
}

type packageLib struct{}

func (packageLib) CompileOptions() []cel.EnvOption {
	reg, err := types.NewRegistry(
		&bundlev1.KV{},
	)
	if err != nil {
		panic(err)
	}

	return []cel.EnvOption{
		cel.Variable("p", harpPackageObjectType),
		cel.Function(
			"match_label",
			cel.MemberOverload("package_match_label_string", []*cel.Type{harpPackageObjectType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(celPackageMatchLabel),
			),
			cel.MemberOverload("package_match_label_string_string", []*cel.Type{harpPackageObjectType, cel.StringType, cel.StringType}, cel.BoolType,
				cel.FunctionBinding(celPackageMatchLabelValue),
			),
		),
		cel.Function(
			"match_annotation",
			cel.MemberOverload("package_match_annotation_string", []*cel.Type{harpPackageObjectType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(celPackageMatchAnnotation),
			),
			cel.MemberOverload("package_match_annotation_string_string", []*cel.Type{harpPackageObjectType, cel.StringType, cel.StringType}, cel.BoolType,
				cel.FunctionBinding(celPackageMatchAnnotationValue),
			),
		),
		cel.Function(
			"match_path",
			cel.MemberOverload("package_match_path_string", []*cel.Type{harpPackageObjectType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(celPackageMatchPath),
			),
		),
		cel.Function(
			"match_secret",
			cel.MemberOverload("package_match_secret_string", []*cel.Type{harpPackageObjectType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(celPackageMatchSecret),
			),
		),
		cel.Function(
			"has_secret",
			cel.MemberOverload("package_has_secret_string", []*cel.Type{harpPackageObjectType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(celPackageHasSecret),
			),
		),
		cel.Function(
			"without_secret",
			cel.MemberOverload("package_without_secret_string", []*cel.Type{harpPackageObjectType}, cel.BoolType,
				cel.UnaryBinding(celPackageWithoutSecret),
			),
		),
		cel.Function(
			"has_all_secrets",
			cel.MemberOverload("package_has_all_secrets_string", []*cel.Type{harpPackageObjectType, cel.ListType(cel.StringType)}, cel.BoolType,
				cel.BinaryBinding(celPackageHasAllSecrets),
			),
		),
		cel.Function(
			"is_cso_compliant",
			cel.MemberOverload("package_is_cso_compliant", []*cel.Type{harpPackageObjectType}, cel.BoolType,
				cel.UnaryBinding(celPackageIsCSOCompliant),
			),
		),
		cel.Function(
			"secret",
			cel.MemberOverload("package_secret_string", []*cel.Type{harpPackageObjectType, cel.StringType}, harpKVObjectType,
				cel.BinaryBinding(celPackageGetSecret(reg)),
			),
		),
	}
}

func (packageLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

// -----------------------------------------------------------------------------

func celPackageMatchLabel(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	patternTyped, ok := rhs.(types.String)
	if !ok {
		return types.Bool(false)
	}

	pattern, ok := patternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	m := glob.MustCompile(pattern)
	for k := range p.Labels {
		if m.Match(k) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageMatchLabelValue(values ...ref.Val) ref.Val {
	if len(values) != 3 {
		return types.Bool(false)
	}

	lhs := values[0]
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	keyPatternTyped, ok := values[1].(types.String)
	if !ok {
		return types.Bool(false)
	}

	keyPattern, ok := keyPatternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	valuePatternTyped, ok := values[2].(types.String)
	if !ok {
		return types.Bool(false)
	}

	valuePattern, ok := valuePatternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	km := glob.MustCompile(keyPattern)
	vm := glob.MustCompile(valuePattern)

	for k, v := range p.Labels {
		if km.Match(k) && vm.Match(v) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageMatchAnnotation(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	patternTyped, ok := rhs.(types.String)
	if !ok {
		return types.Bool(false)
	}

	pattern, ok := patternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	m := glob.MustCompile(pattern)
	for k := range p.Annotations {
		if m.Match(k) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageMatchAnnotationValue(values ...ref.Val) ref.Val {
	if len(values) != 3 {
		return types.Bool(false)
	}

	lhs := values[0]
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	keyPatternTyped, ok := values[1].(types.String)
	if !ok {
		return types.Bool(false)
	}

	keyPattern, ok := keyPatternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	valuePatternTyped, ok := values[2].(types.String)
	if !ok {
		return types.Bool(false)
	}

	valuePattern, ok := valuePatternTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	km := glob.MustCompile(keyPattern)
	vm := glob.MustCompile(valuePattern)

	for k, v := range p.Annotations {
		if km.Match(k) && vm.Match(v) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageMatchPath(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	pathTyped, ok := rhs.(types.String)
	if !ok {
		return types.Bool(false)
	}

	path, ok := pathTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	return types.Bool(glob.MustCompile(path).Match(p.Name))
}

func celPackageMatchSecret(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	secretTyped, ok := rhs.(types.String)
	if !ok {
		return types.Bool(false)
	}

	secretName, ok := secretTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	// No secret data
	if p.Secrets == nil || p.Secrets.Data == nil || len(p.Secrets.Data) == 0 {
		return types.Bool(false)
	}

	m := glob.MustCompile(secretName)

	// Look for secret name
	for _, s := range p.Secrets.Data {
		if m.Match(s.Key) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageHasSecret(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	secretTyped, ok := rhs.(types.String)
	if !ok {
		return types.Bool(false)
	}

	secretName, ok := secretTyped.Value().(string)
	if !ok {
		return types.Bool(false)
	}

	// No secret data
	if p.Secrets == nil || p.Secrets.Data == nil || len(p.Secrets.Data) == 0 {
		return types.Bool(false)
	}

	// Look for secret name
	for _, k := range p.Secrets.Data {
		if strings.EqualFold(k.Key, secretName) {
			return types.Bool(true)
		}
	}

	return types.Bool(false)
}

func celPackageWithoutSecret(value ref.Val) ref.Val {
	x, _ := value.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	// Look for secret name
	if p.Secrets == nil || len(p.Secrets.Data) == 0 {
		return types.Bool(true)
	}

	return types.Bool(false)
}

func celPackageHasAllSecrets(lhs, rhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	secretsTyped, _ := rhs.ConvertToNative(reflect.TypeOf([]string{}))
	secretNames, ok := secretsTyped.([]string)
	if !ok {
		return types.Bool(false)
	}

	// No secret data
	if p.Secrets == nil || p.Secrets.Data == nil || len(p.Secrets.Data) == 0 {
		return types.Bool(false)
	}

	sa := htypes.StringArray(secretNames)

	secretMap := map[string]*bundlev1.KV{}
	for _, k := range p.Secrets.Data {
		if !sa.Contains(k.Key) {
			return types.Bool(false)
		}
		secretMap[k.Key] = k
	}

	// Look for secret name
	for _, k := range secretNames {
		if _, ok := secretMap[k]; !ok {
			return types.Bool(false)
		}
	}

	return types.Bool(true)
}

func celPackageIsCSOCompliant(lhs ref.Val) ref.Val {
	x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
	p, ok := x.(*bundlev1.Package)
	if !ok {
		return types.Bool(false)
	}

	if err := csov1.Validate(p.Name); err != nil {
		return types.Bool(false)
	}

	return types.Bool(true)
}

func celPackageGetSecret(reg ref.TypeAdapter) func(lhs, rhs ref.Val) ref.Val {
	return func(lhs, rhs ref.Val) ref.Val {
		x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.Package{}))
		p, ok := x.(*bundlev1.Package)
		if !ok {
			return types.Bool(false)
		}

		secretTyped, ok := rhs.(types.String)
		if !ok {
			return types.Bool(false)
		}

		secretName, ok := secretTyped.Value().(string)
		if !ok {
			return types.Bool(false)
		}

		// No secret data
		if p.Secrets == nil || p.Secrets.Data == nil || len(p.Secrets.Data) == 0 {
			return types.Bool(false)
		}

		// Look for secret name
		for _, k := range p.Secrets.Data {
			if strings.EqualFold(k.Key, secretName) {
				return reg.NativeToValue(k)
			}
		}

		return nil
	}
}
