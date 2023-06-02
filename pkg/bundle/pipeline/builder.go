// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package pipeline

import "io"

// Options defines default options.
type Options struct {
	input         io.Reader
	output        io.Writer
	disableOutput bool
	fpf           FileProcessorFunc
	ppf           PackageProcessorFunc
	cpf           ChainProcessorFunc
	kpf           KVProcessorFunc
}

// Option represents option function.
type Option func(*Options)

// InputReader defines the input reader used to retrieve the bundle content.
func InputReader(value io.Reader) Option {
	return func(opts *Options) {
		opts.input = value
	}
}

// OutputWriter defines where the bundle will be written after process execution.
func OutputWriter(value io.Writer) Option {
	return func(opts *Options) {
		opts.output = value
	}
}

// OutputDisabled assign the value to disableOutput option.
func OutputDisabled() Option {
	return func(opts *Options) {
		opts.disableOutput = true
	}
}

// FileProcessor assign the file object processor.
func FileProcessor(f FileProcessorFunc) Option {
	return func(opts *Options) {
		opts.fpf = f
	}
}

// PackageProcessor assign the package object processor.
func PackageProcessor(f PackageProcessorFunc) Option {
	return func(opts *Options) {
		opts.ppf = f
	}
}

// ChainProcessor assign the chain object processor.
func ChainProcessor(f ChainProcessorFunc) Option {
	return func(opts *Options) {
		opts.cpf = f
	}
}

// KVProcessor assign the KV object processor.
func KVProcessor(f KVProcessorFunc) Option {
	return func(opts *Options) {
		opts.kpf = f
	}
}
