// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package flatmap

import (
	"fmt"
	"path"
	"reflect"

	"zntr.io/harp/v2/pkg/bundle"
)

// -----------------------------------------------------------------------------

// Flatten takes a structure and turns into a flat map[string]string.
func Flatten(thing map[string]interface{}) map[string]bundle.KV {
	result := make(map[string]string)

	// Flatten recursively the map
	for k, raw := range thing {
		flatten(result, k, reflect.ValueOf(raw))
	}

	// Unpack leaf as secrets
	jsonMap := map[string]bundle.KV{}
	for k, v := range result {
		// Get last element as secret name
		packageName, secretName := path.Split(k)

		// Remove trailing path separator
		packageName = path.Clean(packageName)

		// Check if package already is registered
		p, ok := jsonMap[packageName]
		if !ok {
			p = bundle.KV{}
		}

		// Assign secret
		p[secretName] = v

		// Re-assign to map
		jsonMap[packageName] = p
	}

	// Return json map
	return jsonMap
}

func flatten(result map[string]string, prefix string, v reflect.Value) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			result[prefix] = "true"
		} else {
			result[prefix] = "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result[prefix] = fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result[prefix] = fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		result[prefix] = fmt.Sprintf("%f", v.Float())
	case reflect.Map:
		flattenMap(result, prefix, v)
	case reflect.Slice, reflect.Array:
		flattenSlice(result, prefix, v)
	case reflect.String:
		result[prefix] = v.String()
	case reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func, reflect.Interface:
		// ignore
	case reflect.Invalid, reflect.Ptr, reflect.Struct, reflect.Uintptr, reflect.UnsafePointer:
		// ignore
	default:
		panic(fmt.Sprintf("Unknown: %s", v))
	}
}

func flattenMap(result map[string]string, prefix string, v reflect.Value) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}

		flatten(result, fmt.Sprintf("%s/%s", prefix, k.String()), v.MapIndex(k))
	}
}

func flattenSlice(result map[string]string, prefix string, v reflect.Value) {
	prefix += "/"

	result[prefix+"#"] = fmt.Sprintf("%d", v.Len())
	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i))
	}
}
