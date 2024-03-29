// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	bundlev1 "zntr.io/harp/v2/api/gen/go/harp/bundle/v1"
	"zntr.io/harp/v2/pkg/bundle/secret"
	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestDumpTask_Run(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		PathOnly        bool
		DataOnly        bool
		MetadataOnly    bool
		JMESPathFilter  string
		IgnoreTemplate  bool
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
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("non-existent.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "containerReader not a bundle",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: true,
		},
		{
			name: "outputWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
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
				OutputWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
			},
			wantErr: true,
		},
		{
			name: "invalid JMES Filter",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				JMESPathFilter:  ".",
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid - path only",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				PathOnly:        true,
			},
			wantErr: false,
		},
		{
			name: "valid - data only",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				DataOnly:        true,
			},
			wantErr: false,
		},
		{
			name: "valid - metadata only",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				MetadataOnly:    true,
			},
			wantErr: false,
		},
		{
			name: "valid - with JMES Filter",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
				JMESPathFilter:  "merkleTreeRoot",
			},
			wantErr: false,
		},
		{
			name: "valid",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				OutputWriter:    cmdutil.DiscardWriter(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DumpTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				PathOnly:        tt.fields.PathOnly,
				DataOnly:        tt.fields.DataOnly,
				MetadataOnly:    tt.fields.MetadataOnly,
				JMESPathFilter:  tt.fields.JMESPathFilter,
				IgnoreTemplate:  tt.fields.IgnoreTemplate,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DumpTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDumpTask_dumpData_NilWriter(t *testing.T) {
	tr := &DumpTask{}
	err := tr.dumpData(nil, nil)
	assert.Error(t, err)
}

func TestDumpTask_dumpData(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		PathOnly        bool
		DataOnly        bool
		MetadataOnly    bool
		JMESPathFilter  string
	}
	type args struct {
		b *bundlev1.Bundle
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "empty bundle",
			args: args{
				b: &bundlev1.Bundle{},
			},
			wantWriter: `{}` + "\n",
			wantErr:    false,
		},
		{
			name: "bundle with secrets",
			args: args{
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secret/package",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "test",
										Value: secret.MustPack("value"),
									},
								},
							},
						},
					},
				},
			},
			wantWriter: `{"secret/package":{"test":"value"}}` + "\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DumpTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				PathOnly:        tt.fields.PathOnly,
				DataOnly:        tt.fields.DataOnly,
				MetadataOnly:    tt.fields.MetadataOnly,
				JMESPathFilter:  tt.fields.JMESPathFilter,
			}
			writer := &bytes.Buffer{}
			if err := tr.dumpData(writer, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("DumpTask.dumpData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("DumpTask.dumpData() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestDumpTask_dumpMetadata_NilWriter(t *testing.T) {
	tr := &DumpTask{}
	err := tr.dumpMetadata(nil, nil)
	assert.Error(t, err)
}

func TestDumpTask_dumpMetadata(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		PathOnly        bool
		DataOnly        bool
		MetadataOnly    bool
		JMESPathFilter  string
	}
	type args struct {
		b *bundlev1.Bundle
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "empty bundle",
			args: args{
				b: &bundlev1.Bundle{},
			},
			wantWriter: `{}` + "\n",
			wantErr:    false,
		},
		{
			name: "bundle with secrets only",
			args: args{
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secret/package",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "test",
										Value: secret.MustPack("value"),
									},
								},
							},
						},
					},
				},
			},
			wantWriter: `{"secret/package":{}}` + "\n",
			wantErr:    false,
		},
		{
			name: "bundle with secrets and metadata",
			args: args{
				b: &bundlev1.Bundle{
					Labels: map[string]string{
						"generated": "true",
					},
					Annotations: map[string]string{
						"annotation": "text",
					},
					Packages: []*bundlev1.Package{
						{
							Name: "secret/package",
							Labels: map[string]string{
								"vendor": "true",
							},
							Annotations: map[string]string{
								"harp.zntr.io/v2/testing#Annotation": "test",
							},
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "test",
										Value: secret.MustPack("value"),
									},
								},
							},
						},
					},
				},
			},
			wantWriter: `{"harp.zntr.io/v2/bundle#annotations":{"annotation":"text"},"harp.zntr.io/v2/bundle#labels":{"generated":"true"},"secret/package":{"harp.zntr.io/v2/package#annotations":{"harp.zntr.io/v2/testing#Annotation":"test"},"harp.zntr.io/v2/package#labels":{"vendor":"true"}}}` + "\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DumpTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				PathOnly:        tt.fields.PathOnly,
				DataOnly:        tt.fields.DataOnly,
				MetadataOnly:    tt.fields.MetadataOnly,
				JMESPathFilter:  tt.fields.JMESPathFilter,
			}
			writer := &bytes.Buffer{}
			if err := tr.dumpMetadata(writer, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("DumpTask.dumpMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("DumpTask.dumpMetadata() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestDumpTask_dumpPath_NilWriter(t *testing.T) {
	tr := &DumpTask{}
	err := tr.dumpPath(nil, nil)
	assert.Error(t, err)
}

func TestDumpTask_dumpPath(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		OutputWriter    tasks.WriterProvider
		PathOnly        bool
		DataOnly        bool
		MetadataOnly    bool
		JMESPathFilter  string
	}
	type args struct {
		b *bundlev1.Bundle
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:    "nil",
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "empty bundle",
			args: args{
				b: &bundlev1.Bundle{},
			},
			wantWriter: "",
			wantErr:    false,
		},
		{
			name: "bundle with secrets",
			args: args{
				b: &bundlev1.Bundle{
					Packages: []*bundlev1.Package{
						{
							Name: "secret/package",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "test",
										Value: secret.MustPack("value"),
									},
								},
							},
						},
						{
							Name: "application/security",
							Secrets: &bundlev1.SecretChain{
								Data: []*bundlev1.KV{
									{
										Key:   "test",
										Value: secret.MustPack("value"),
									},
								},
							},
						},
					},
				},
			},
			wantWriter: "application/security\nsecret/package\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DumpTask{
				ContainerReader: tt.fields.ContainerReader,
				OutputWriter:    tt.fields.OutputWriter,
				PathOnly:        tt.fields.PathOnly,
				DataOnly:        tt.fields.DataOnly,
				MetadataOnly:    tt.fields.MetadataOnly,
				JMESPathFilter:  tt.fields.JMESPathFilter,
			}
			writer := &bytes.Buffer{}
			if err := tr.dumpPath(writer, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("DumpTask.dumpPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("DumpTask.dumpPath() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestDumpTask_dumpFilter_NilWriter(t *testing.T) {
	tr := &DumpTask{}
	err := tr.dumpFilter(nil, nil)
	assert.Error(t, err)

	err = tr.dumpFilter(&bytes.Buffer{}, nil)
	assert.Error(t, err)
}
