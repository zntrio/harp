// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

import (
	"encoding/base64"
	"path"
	"strings"

	"github.com/gobwas/glob"
	"zntr.io/harp/v2/pkg/template/files"
)

// Files is a map of files that can be accessed from a template.
type Files map[string][]byte

// NewFiles returns an engine file collection from file loader.
func NewFiles(from []*files.BufferedFile) Files {
	fileMap := make(map[string][]byte)
	for _, f := range from {
		fileMap[f.Name] = f.Data
	}

	return fileMap
}

// GetBytes gets a file by path.
//
// The returned data is raw. In a template context, this is identical to calling
// {{index .Files $path}}.
//
// This is intended to be accessed from within a template, so a missed key returns
// an empty []byte.
func (f Files) GetBytes(name string) []byte {
	if v, ok := f[name]; ok {
		return v
	}
	return []byte{}
}

// Get returns a string representation of the given file.
//
// Fetch the contents of a file as a string. It is designed to be called in a
// template.
//
//	{{.Files.Get "foo"}}
func (f Files) Get(name string) string {
	return string(f.GetBytes(name))
}

// Glob takes a glob pattern and returns another files object only containing
// matched  files.
//
// This is designed to be called from a template.
//
// {{ range $name, $content := .Files.Glob("foo/**") }}
// {{ $name }}: |
// {{ .Files.Get($name) | indent 4 }}{{ end }}.
func (f Files) Glob(pattern string) Files {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		g, _ = glob.Compile("**")
	}

	nf := Files{}
	for name, contents := range f {
		if g.Match(name) {
			nf[name] = contents
		}
	}

	return nf
}

// AsConfig returns a Files group and flattens it to a YAML map suitable for
// including in the 'data' section of a Kubernetes ConfigMap definition.
// Duplicate keys will be overwritten, so be aware that your file names
// (regardless of path) should be unique.
//
// The output will not be indented, so you will want to pipe this to the
// 'indent' template function.
//
//	data:
//
// {{ .Files.Glob("config/**").AsConfig() | toYaml | indent 4 }}.
func (f Files) AsConfig() map[string]string {
	if f == nil {
		return nil
	}

	m := make(map[string]string)

	// Explicitly convert to strings, and file names
	for k, v := range f {
		m[path.Base(k)] = string(v)
	}

	return m
}

// AsSecrets returns the base64-encoded value of a Files object suitable for
// including in the 'data' section of a Kubernetes Secret definition.
// Duplicate keys will be overwritten, so be aware that your file names
// (regardless of path) should be unique.
//
// The output will not be indented, so you will want to pipe this to the
// 'indent' template function.
//
//	data:
//
// {{ .Files.Glob("secrets/*").AsSecrets() | toYaml }}.
func (f Files) AsSecrets() map[string]string {
	if f == nil {
		return nil
	}

	m := make(map[string]string)

	for k, v := range f {
		m[path.Base(k)] = base64.StdEncoding.EncodeToString(v)
	}

	return m
}

// Lines returns each line of a named file (split by "\n") as a slice, so it can
// be ranged over in your templates.
//
// This is designed to be called from a template.
//
// {{ range .Files.Lines "foo/bar.html" }}
// {{ . }}{{ end }}.
func (f Files) Lines(filePath string) []string {
	if f == nil || f[filePath] == nil {
		return []string{}
	}

	return strings.Split(string(f[filePath]), "\n")
}
