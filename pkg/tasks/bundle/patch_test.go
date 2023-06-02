// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"errors"
	"io"
	"testing"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestPatchTask_Run(t *testing.T) {
	type fields struct {
		PatchReader     tasks.ReaderProvider
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		Values          map[string]interface{}
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
			name: "nil patchReader",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     nil,
			},
			wantErr: true,
		},
		{
			name: "nil outputWriter",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter:    nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("non-existent.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "containerReader - not a bundle",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "patchReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("non-existent.yaml"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "patchReader - not a yaml",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
			},
			wantErr: true,
		},
		{
			name: "outputWriter closed",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/path-cleaner.yaml"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: false,
		},
		{
			name: "empty patch",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				PatchReader:     cmdutil.FileReader("../../../test/fixtures/patch/valid/empty.yaml"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &PatchTask{
				PatchReader:     tt.fields.PatchReader,
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				Values:          tt.fields.Values,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("PatchTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
