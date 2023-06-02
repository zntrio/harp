// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

import (
	"bytes"
	"fmt"
	"text/template"
)

// -----------------------------------------------------------------------------

// Render compile and assemble attribute template to merge with values.
func Render(input string, data interface{}) (content string, err error) {
	// Check argument
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("template rendering failed: %v", r)
		}
	}()

	// Prepare the template
	t, err := template.New("root").
		Funcs(FuncMap(nil)).
		Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to compile attribute template %q: %w", input, err)
	}

	// Fail on missing key
	t.Option("missingkey=error")

	// Merge with values
	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", fmt.Errorf("unable to merge data with template %q: %w", input, err)
	}

	// No error
	return out.String(), nil
}

// RenderContext compile and assemble attribute template to merge with values.
func RenderContext(templateContext Context, input string) (string, error) {
	return RenderContextWithData(templateContext, input, nil)
}

// RenderContextWithData compile and assemble attribute template to merge with values.
func RenderContextWithData(templateContext Context, input string, data interface{}) (content string, err error) {
	// Check argument
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("template rendering failed: %v", r)
		}
	}()

	// Retrieve delimiters
	leftDelim, rightDelim := templateContext.Delims()

	// Prepare the template
	t, err := template.New(templateContext.Name()).
		Delims(leftDelim, rightDelim).
		Funcs(FuncMap(templateContext.SecretReaders())).
		Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to compile attribute template %q: %w", input, err)
	}

	// Check strict mode
	if templateContext.StrictMode() {
		// Fail on missing key
		t.Option("missingkey=error")
	} else {
		// Not that zero will attempt to add default values for types it knows,
		// but will still emit <no value> for others. We mitigate that later.
		t.Option("missingkey=zero")
	}

	// Merge with values
	var out bytes.Buffer
	if err := t.Execute(&out, map[string]interface{}{
		"Data":   data,
		"Values": templateContext.Values(),
		"Files":  templateContext.Files(),
	}); err != nil {
		return "", fmt.Errorf("unable to merge values with template %q: %w", input, err)
	}

	// No error
	return out.String(), nil
}
