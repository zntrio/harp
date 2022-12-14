// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bundle

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/zntrio/harp/v2/pkg/sdk/cmdutil"
	"github.com/zntrio/harp/v2/pkg/tasks"
)

func TestDiffTask_Run(t *testing.T) {
	type fields struct {
		SourceReader      tasks.ReaderProvider
		DestinationReader tasks.ReaderProvider
		OutputWriter      tasks.WriterProvider
		GeneratePatch     bool
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
			name: "nil destinationReader",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: nil,
			},
			wantErr: true,
		},
		{
			name: "nil outputWriter",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      nil,
			},
			wantErr: true,
		},
		{
			name: "sourceReader error",
			fields: fields{
				SourceReader:      cmdutil.FileReader("non-existent.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: true,
		},
		{
			name: "sourceReader not a bundle",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: true,
		},
		{
			name: "destinationReader error",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("non-existent.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: true,
		},
		{
			name: "sourceReader not a bundle",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
			},
			wantErr: true,
		},
		{
			name: "outputWriter closed",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
			},
			wantErr: true,
		},
		{
			name: "self-diff - with patch",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     true,
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "self-diff",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: false,
		},
		{
			name: "bundle diff",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/empty.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     false,
			},
			wantErr: false,
		},
		{
			name: "bundle diff - with patch",
			fields: fields{
				SourceReader:      cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				DestinationReader: cmdutil.FileReader("../../../test/fixtures/bundles/empty.bundle"),
				OutputWriter:      cmdutil.DiscardWriter(),
				GeneratePatch:     true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DiffTask{
				SourceReader:      tt.fields.SourceReader,
				DestinationReader: tt.fields.DestinationReader,
				OutputWriter:      tt.fields.OutputWriter,
				GeneratePatch:     tt.fields.GeneratePatch,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DiffTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
