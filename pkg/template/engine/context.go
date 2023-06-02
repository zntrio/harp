// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package engine

// Context describes engine rendering context contract.
type Context interface {
	Name() string
	StrictMode() bool
	Delims() (string, string)
	SecretReaders() []SecretReaderFunc
	Values() Values
	Files() Files
}

// -----------------------------------------------------------------------------

// ContextOption defines context functional builder function.
type ContextOption func(*context)

// WithName sets the template name.
func WithName(value string) ContextOption {
	return func(ctx *context) {
		ctx.name = value
	}
}

// WithStrictMode enable or disable strict rendering mode.
func WithStrictMode(value bool) ContextOption {
	return func(ctx *context) {
		ctx.strictMode = value
	}
}

// WithDelims defines used delimiters for rendering engine.
func WithDelims(left, right string) ContextOption {
	return func(ctx *context) {
		ctx.delimLeft = left
		ctx.delimRight = right
	}
}

// WithSecretReaders defines secret resolver functions used by `secret` template
// function.
func WithSecretReaders(values ...SecretReaderFunc) ContextOption {
	return func(ctx *context) {
		if len(values) > 0 {
			ctx.secretReaders = values
		}
	}
}

// WithValues defines template values injected via CLI.
func WithValues(values Values) ContextOption {
	return func(ctx *context) {
		ctx.values = values
	}
}

// WithFiles defines file collection.
func WithFiles(files Files) ContextOption {
	return func(ctx *context) {
		ctx.files = files
	}
}

// NewContext returns a template rendering context.
func NewContext(opts ...ContextOption) Context {
	defaultContext := &context{
		delimLeft:     "{{",
		delimRight:    "}}",
		name:          "root",
		secretReaders: []SecretReaderFunc{},
		strictMode:    true,
	}

	// Apply functions
	for _, opt := range opts {
		opt(defaultContext)
	}

	// Return modified context
	return defaultContext
}

// -----------------------------------------------------------------------------

// Context describes rendering context.
type context struct {
	name          string
	strictMode    bool
	delimLeft     string
	delimRight    string
	secretReaders []SecretReaderFunc
	values        Values
	files         Files
}

// Name returns template name.
func (ctx *context) Name() string {
	return ctx.name
}

// StrictMode retruns strict mode status of template engine.
func (ctx *context) StrictMode() bool {
	return ctx.strictMode
}

// Delims returns left and right delimiters used to compile the template.
func (ctx *context) Delims() (left, right string) {
	return ctx.delimLeft, ctx.delimRight
}

// SecretReaders returns secret reader function called by `secret` template function.
func (ctx *context) SecretReaders() []SecretReaderFunc {
	return ctx.secretReaders
}

// Values returns binded values from rendering context.
func (ctx *context) Values() Values {
	return ctx.values
}

// Files returns binded files from rendering context.
func (ctx *context) Files() Files {
	return ctx.files
}
