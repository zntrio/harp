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
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"
	"zntr.io/harp/v2/pkg/sdk/value/identity"
	"zntr.io/harp/v2/pkg/sdk/value/mock"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestEncryptTask_Run(t *testing.T) {
	type fields struct {
		ContainerReader   tasks.ReaderProvider
		OutputWriter      tasks.WriterProvider
		BundleTransformer value.Transformer
		TransformerMap    map[string]value.Transformer
		SkipUnresolved    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "nil",
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
			name: "nil transformers",
			fields: fields{
				ContainerReader:   cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				BundleTransformer: nil,
				TransformerMap:    nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader:   cmdutil.FileReader("non-existent.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				BundleTransformer: identity.Transformer(),
			},
			wantErr: true,
		},
		{
			name: "containerReader - not a bundle",
			fields: fields{
				ContainerReader:   cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				OutputWriter:      cmdutil.DiscardWriter(),
				BundleTransformer: identity.Transformer(),
			},
			wantErr: true,
		},
		{
			name: "bundle transformer error",
			fields: fields{
				ContainerReader:   cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				BundleTransformer: mock.Transformer(errors.New("test")),
			},
			wantErr: true,
		},
		{
			name: "empty annotation transformer map",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				TransformerMap:  map[string]value.Transformer{},
			},
			wantErr: true,
		},
		{
			name: "annotation transformer error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				TransformerMap: map[string]value.Transformer{
					"test": mock.Transformer(errors.New("test")),
				},
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
				BundleTransformer: encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
			},
			wantErr: true,
		},
		{
			name: "outputWriter closed",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
				BundleTransformer: encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader:   cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				BundleTransformer: encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
			},
			wantErr: false,
		},
		{
			name: "valid - unused key alias",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				TransformerMap: map[string]value.Transformer{
					"not-used": encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
				},
				SkipUnresolved: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &EncryptTask{
				ContainerReader:   tt.fields.ContainerReader,
				OutputWriter:      tt.fields.OutputWriter,
				BundleTransformer: tt.fields.BundleTransformer,
				TransformerMap:    tt.fields.TransformerMap,
				SkipUnresolved:    tt.fields.SkipUnresolved,
			}
			if err := tr.Run(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("EncryptTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
