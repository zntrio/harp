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
	fuzz "github.com/google/gofuzz"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestSealTask_Run_V1(t *testing.T) {
	pub := "v1.sk.qKXPnUP6-2Bb_4nYnmxOXyCdN4IV3AR5HooB33N3g2E"

	type fields struct {
		ContainerReader          tasks.ReaderProvider
		SealedContainerWriter    tasks.WriterProvider
		OutputWriter             tasks.WriterProvider
		PeerPublicKeys           []string
		DCKDMasterKey            string
		DCKDTarget               string
		JSONOutput               bool
		DisableContainerIdentity bool
		PreSharedKey             *memguard.LockedBuffer
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
			name: "nil sealedContainerWriter",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: nil,
			},
			wantErr: true,
		},
		{
			name: "nil outputWriter",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          nil,
			},
			wantErr: true,
		},
		{
			name: "no public keys",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{},
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("non-existent.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pub},
			},
			wantErr: true,
		},
		{
			name: "containerReader not a bundle",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pub},
			},
			wantErr: true,
		},
		{
			name: "sealedContainerWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
				OutputWriter:   cmdutil.DiscardWriter(),
				PeerPublicKeys: []string{pub},
			},
			wantErr: true,
		},
		{
			name: "sealedContainerWriter closed",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
				OutputWriter:   cmdutil.DiscardWriter(),
				PeerPublicKeys: []string{pub},
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pub},
			},
			wantErr: false,
		},
		{
			name: "valid with psk",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pub},
				PreSharedKey:          memguard.NewBufferFromBytes([]byte("Kw6tb0QWUH3vueG5uCvS6lAnUa00a5-lsM2aqOZk3MFvoDTUUyhjIdb6ZAG7eQt3LJ1QnJQQAZBLVGXQkx33kg")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &SealTask{
				ContainerReader:          tt.fields.ContainerReader,
				SealedContainerWriter:    tt.fields.SealedContainerWriter,
				OutputWriter:             tt.fields.OutputWriter,
				PeerPublicKeys:           tt.fields.PeerPublicKeys,
				DCKDMasterKey:            tt.fields.DCKDMasterKey,
				DCKDTarget:               tt.fields.DCKDTarget,
				JSONOutput:               tt.fields.JSONOutput,
				DisableContainerIdentity: tt.fields.DisableContainerIdentity,
				SealVersion:              1,
				PreSharedKey:             tt.fields.PreSharedKey,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SealTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSealTask_Run_V2(t *testing.T) {
	pk := "v2.sk.A0V1xCxGNtVAE9EVhaKi-pIADhd1in8xV_FI5Y0oHSHLAkew9gDAqiALSd6VgvBCbQ"

	type fields struct {
		ContainerReader          tasks.ReaderProvider
		SealedContainerWriter    tasks.WriterProvider
		OutputWriter             tasks.WriterProvider
		PeerPublicKeys           []string
		DCKDMasterKey            string
		DCKDTarget               string
		JSONOutput               bool
		DisableContainerIdentity bool
		PreSharedKey             *memguard.LockedBuffer
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
			name: "nil sealedContainerWriter",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: nil,
			},
			wantErr: true,
		},
		{
			name: "nil outputWriter",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          nil,
			},
			wantErr: true,
		},
		{
			name: "no public keys",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{},
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("non-existent.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pk},
			},
			wantErr: true,
		},
		{
			name: "containerReader not a bundle",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pk},
			},
			wantErr: true,
		},
		{
			name: "sealedContainerWriter error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: func(ctx context.Context) (io.Writer, error) {
					return nil, errors.New("test")
				},
				OutputWriter:   cmdutil.DiscardWriter(),
				PeerPublicKeys: []string{pk},
			},
			wantErr: true,
		},
		{
			name: "sealedContainerWriter closed",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: func(ctx context.Context) (io.Writer, error) {
					return cmdutil.NewClosedWriter(), nil
				},
				OutputWriter:   cmdutil.DiscardWriter(),
				PeerPublicKeys: []string{pk},
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pk},
			},
			wantErr: false,
		},
		{
			name: "valid with psk",
			fields: fields{
				ContainerReader:       cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				SealedContainerWriter: cmdutil.DiscardWriter(),
				OutputWriter:          cmdutil.DiscardWriter(),
				PeerPublicKeys:        []string{pk},
				PreSharedKey:          memguard.NewBufferFromBytes([]byte("Kw6tb0QWUH3vueG5uCvS6lAnUa00a5-lsM2aqOZk3MFvoDTUUyhjIdb6ZAG7eQt3LJ1QnJQQAZBLVGXQkx33kg")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &SealTask{
				ContainerReader:          tt.fields.ContainerReader,
				SealedContainerWriter:    tt.fields.SealedContainerWriter,
				OutputWriter:             tt.fields.OutputWriter,
				PeerPublicKeys:           tt.fields.PeerPublicKeys,
				DCKDMasterKey:            tt.fields.DCKDMasterKey,
				DCKDTarget:               tt.fields.DCKDTarget,
				JSONOutput:               tt.fields.JSONOutput,
				DisableContainerIdentity: tt.fields.DisableContainerIdentity,
				SealVersion:              2,
				PreSharedKey:             tt.fields.PreSharedKey,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SealTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSealTask_Fuzz(t *testing.T) {
	tsk := &SealTask{
		ContainerReader:          cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
		SealedContainerWriter:    cmdutil.DiscardWriter(),
		OutputWriter:             cmdutil.DiscardWriter(),
		PeerPublicKeys:           []string{},
		DisableContainerIdentity: true,
	}

	// Making sure the function never panics
	for i := 0; i < 50; i++ {
		f := fuzz.New()

		// Prepare arguments
		f.Fuzz(&tsk.PeerPublicKeys)

		// Execute
		tsk.Run(context.Background())
	}
}
