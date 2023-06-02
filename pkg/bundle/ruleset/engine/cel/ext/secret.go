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
	"encoding/json"
	"fmt"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/secret"
)

// Secrets exported secret operations.
func Secrets() cel.EnvOption {
	return cel.Lib(secretLib{})
}

type secretLib struct{}

func (secretLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function(
			"is_base64",
			cel.MemberOverload("kv_is_base64", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(is.Base64)),
			),
		),
		cel.Function(
			"is_required",
			cel.MemberOverload("kv_is_required", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(validation.Required)),
			),
		),
		cel.Function(
			"is_url",
			cel.MemberOverload("kv_is_url", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(is.URL)),
			),
		),
		cel.Function(
			"is_uuid",
			cel.MemberOverload("kv_is_uuid", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(is.UUID)),
			),
		),
		cel.Function(
			"is_email",
			cel.MemberOverload("kv_is_email", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(is.EmailFormat)),
			),
		),
		cel.Function(
			"is_json",
			cel.MemberOverload("kv_is_json", []*cel.Type{harpKVObjectType}, cel.BoolType,
				cel.UnaryBinding(celValidatorBuilder(&jsonValidator{})),
			),
		),
	}
}

func (secretLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

// -----------------------------------------------------------------------------

func celValidatorBuilder(rules ...validation.Rule) func(ref.Val) ref.Val {
	return func(lhs ref.Val) ref.Val {
		x, _ := lhs.ConvertToNative(reflect.TypeOf(&bundlev1.KV{}))
		p, ok := x.(*bundlev1.KV)
		if !ok {
			return types.Bool(false)
		}

		var out string
		if err := secret.Unpack(p.Value, &out); err != nil {
			return types.Bool(false)
		}

		if err := validation.Validate(out, rules...); err != nil {
			return types.Bool(false)
		}

		return types.Bool(true)
	}
}

// -----------------------------------------------------------------------------

var _ validation.Rule = (*jsonValidator)(nil)

type jsonValidator struct{}

func (v *jsonValidator) Validate(in interface{}) error {
	// Process input
	if data, ok := in.([]byte); ok {
		if !json.Valid(data) {
			return fmt.Errorf("invalid JSON payload")
		}
	}

	return fmt.Errorf("unable to validate JSON for %T type", in)
}
