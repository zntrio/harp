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

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/identity"
	"zntr.io/harp/v2/pkg/sdk/value/mock"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestIdentityTask_Run(t *testing.T) {
	type fields struct {
		OutputWriter tasks.WriterProvider
		Description  string
		Transformer  value.Transformer
		Version      IdentityVersion
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
			name: "nil outputWriter",
			fields: fields{
				OutputWriter: nil,
			},
			wantErr: true,
		},
		{
			name: "nil transformer",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Transformer:  nil,
			},
			wantErr: true,
		},
		{
			name: "blank description",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Transformer:  identity.Transformer(),
				Description:  "",
			},
			wantErr: true,
		},
		{
			name: "transformer error",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Transformer:  mock.Transformer(errors.New("test")),
				Description:  "test",
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
				Description: "test",
				Transformer: identity.Transformer(),
			},
			wantErr: true,
		},
		{
			name: "outputWriter closed",
			fields: fields{
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
				Description: "test",
				Transformer: identity.Transformer(),
			},
			wantErr: true,
		},
		{
			name: "version unspecified",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Description:  "test",
				Transformer:  identity.Transformer(),
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid - v1",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Description:  "test",
				Transformer:  identity.Transformer(),
				Version:      LegacyIdentity,
			},
			wantErr: false,
		},
		{
			name: "valid - v2",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Description:  "test",
				Transformer:  identity.Transformer(),
				Version:      ModernIdentity,
			},
			wantErr: false,
		},
		{
			name: "valid - v3",
			fields: fields{
				OutputWriter: cmdutil.DiscardWriter(),
				Description:  "test",
				Transformer:  identity.Transformer(),
				Version:      NISTIdentity,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &IdentityTask{
				OutputWriter: tt.fields.OutputWriter,
				Description:  tt.fields.Description,
				Transformer:  tt.fields.Transformer,
				Version:      tt.fields.Version,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("IdentityTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
