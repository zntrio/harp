// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/awnumar/memguard"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestUnsealTask_Run(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		ContainerKey    *memguard.LockedBuffer
		PreSharedKey    *memguard.LockedBuffer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		{
			name: "nil containerReader",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
			},
			wantErr: true,
		},
		{
			name: "nil outputWriter",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    nil,
			},
			wantErr: true,
		},
		{
			name: "nil containerKey",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("non-existent.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBuffer(32),
			},
			wantErr: true,
		},
		{
			name: "containerReader not a bundle",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBuffer(32),
			},
			wantErr: true,
		},
		{
			name: "invalid container key",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBuffer(32),
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v1.sealed"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
				ContainerKey: memguard.NewBufferFromBytes([]byte("v1.ck.MiVGh4KOmdzZbej17BZGChkCPZ9uK9uBWdPNU0GlBNg")),
			},
			wantErr: true,
		},
		{
			name: "outputWriter closed",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v1.sealed"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
				ContainerKey: memguard.NewBufferFromBytes([]byte("v1.ck.MiVGh4KOmdzZbej17BZGChkCPZ9uK9uBWdPNU0GlBNg")),
			},
			wantErr: true,
		},
		{
			name: "v2 without prefix",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v2.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v2.ck.dAYx4CeTMRGKfpFHA7Q926qMz8imo1VJIToMw9uvH7HfPJTRpLUSMUS07JAdV-1R")),
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid - v1",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v1.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v1.ck.MiVGh4KOmdzZbej17BZGChkCPZ9uK9uBWdPNU0GlBNg")),
			},
			wantErr: false,
		},
		{
			name: "valid - v1 - with identity recovery key",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v1.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v1.ck.IO6bCjACnqsCP0ahT--CVBhryzhe-ZFroVzn5Dx3D0U")),
			},
			wantErr: false,
		},
		{
			name: "valid - v1 - with identity recovery key with prefix",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v1.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v1.ck.IO6bCjACnqsCP0ahT--CVBhryzhe-ZFroVzn5Dx3D0U")),
			},
			wantErr: false,
		},
		{
			name: "valid - v2",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v2.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v2.ck.CLMEUoY-EgvMGKCcKeByPdJjQDod6fqTnqvxtD_Z0_SX4PMITu_emttDL91z_61D")),
			},
			wantErr: false,
		},
		{
			name: "valid - v2 - with identity recovery key",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.v2.sealed"),
				OutputWriter:    cmdutil.DiscardWriter(),
				ContainerKey:    memguard.NewBufferFromBytes([]byte("v2.ck.8DwD8D-TUB9w-NzXBXySz4PkAIrWUc09TOJKdJ495MJ-AJ2lvDlj1Pnw1rSUAwVg")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &UnsealTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				ContainerKey:    tt.fields.ContainerKey,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("UnsealTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
