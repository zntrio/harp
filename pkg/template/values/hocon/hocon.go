// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package hocon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-akka/configuration"
	"github.com/go-akka/configuration/hocon"
	"go.uber.org/zap"
	"zntr.io/harp/v2/pkg/sdk/log"
)

// Parser is a HOCON parser.
type Parser struct{}

// Unmarshal unmarshals HOCON files.
func (i *Parser) Unmarshal(p []byte, v interface{}) error {
	// Parse HOCON configuration
	rootCfg := configuration.ParseString(string(p), hoconIncludeCallback).Root()

	// Visit config tree
	res := visitNode(rootCfg)

	// Encode as json
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(res); err != nil {
		return fmt.Errorf("unable to encode HOCON map to JSON: %w", err)
	}

	// Decode JSON
	if err := json.Unmarshal(buf.Bytes(), v); err != nil {
		return fmt.Errorf("unable to decode json object as struct: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------

func visitNode(node *hocon.HoconValue) interface{} {
	if node.IsArray() {
		nodes := node.GetArray()

		res := make([]interface{}, len(nodes))
		for i, n := range nodes {
			res[i] = visitNode(n)
		}

		return res
	}

	if node.IsObject() {
		obj := node.GetObject()

		res := map[string]interface{}{}
		keys := obj.GetKeys()
		for _, k := range keys {
			res[k] = visitNode(obj.GetKey(k))
		}

		return res
	}

	if node.IsString() {
		return node.GetString()
	}

	if node.IsEmpty() {
		return nil
	}

	return nil
}

func hoconIncludeCallback(filename string) *hocon.HoconRoot {
	files, err := filepath.Glob(filename)
	switch {
	case err != nil:
		log.Bg().Error("hocon: unable to load file glob", zap.Error(err), zap.String("filename", filename))
		return nil
	case len(files) == 0:
		log.Bg().Warn("hocon: unable to load file %s", zap.String("filename", filename))
		return hocon.Parse("", nil)
	default:
		root := hocon.Parse("", nil)
		for _, f := range files {
			data, err := os.ReadFile(filepath.Clean(f))
			if err != nil {
				log.Bg().Error("hocon: unable to load file glob", zap.Error(err))
				return nil
			}

			node := hocon.Parse(string(data), hoconIncludeCallback)
			if node != nil {
				root.Value().GetObject().Merge(node.Value().GetObject())
				// merge substitutions
				subs := make([]*hocon.HoconSubstitution, 0)
				subs = append(subs, root.Substitutions()...)
				subs = append(subs, node.Substitutions()...)
				root = hocon.NewHoconRoot(root.Value(), subs...)
			}
		}
		return root
	}
}
