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

	// Import for tests.
	_ "zntr.io/harp/v2/pkg/sdk/value/encryption/aead"
	"zntr.io/harp/v2/pkg/sdk/value/identity"
	"zntr.io/harp/v2/pkg/sdk/value/mock"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestDecryptTask_Run(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		Transformers    []value.Transformer
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
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter:    nil,
			},
			wantErr: true,
		},
		{
			name: "nil transformer",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers:    nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("non-existent.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers: []value.Transformer{
					identity.Transformer(),
				},
			},
			wantErr: true,
		},
		{
			name: "containerReader - not a bundle",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers: []value.Transformer{
					identity.Transformer(),
				},
			},
			wantErr: true,
		},
		{
			name: "transformer error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers: []value.Transformer{
					mock.Transformer(errors.New("test")),
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
				Transformers: []value.Transformer{
					encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
				},
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
				Transformers: []value.Transformer{
					encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
				},
			},
			wantErr: true,
		},
		{
			name: "no valid key provided",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers: []value.Transformer{
					encryption.Must(encryption.FromKey("aes-gcm:h_0H0n0w0c0c1bw7_orRoA==")),
				},
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.aes-gcm.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				Transformers: []value.Transformer{
					encryption.Must(encryption.FromKey("aes-gcm:5OSpiJUr_XS2M1_vvTBeGg==")),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DecryptTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				Transformers:    tt.fields.Transformers,
			}
			if err := tr.Run(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("DecryptTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
