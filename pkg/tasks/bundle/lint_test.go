// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package bundle

import (
	"context"
	"testing"

	"zntr.io/harp/v2/pkg/sdk/cmdutil"
	"zntr.io/harp/v2/pkg/tasks"
)

func TestLintTask_Run(t *testing.T) {
	type fields struct {
		ContainerReader tasks.ReaderProvider
		RuleSetReader   tasks.ReaderProvider
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
			name: "nil ruleSetReader",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				RuleSetReader:   nil,
			},
			wantErr: true,
		},
		{
			name: "containerReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("non-existent.bundle"),
				RuleSetReader:   cmdutil.FileReader("../../../test/fixtures/ruleset/valid/cso.yaml"),
			},
			wantErr: true,
		},
		{
			name: "containerReader - not a bundle",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.json"),
				RuleSetReader:   cmdutil.FileReader("../../../test/fixtures/ruleset/valid/cso.yaml"),
			},
			wantErr: true,
		},
		{
			name: "ruleSetReader error",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				RuleSetReader:   cmdutil.FileReader("non-existent.yaml"),
			},
			wantErr: true,
		},
		{
			name: "containerReader - not a yaml",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				RuleSetReader:   cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
			},
			wantErr: true,
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				RuleSetReader:   cmdutil.FileReader("../../../test/fixtures/ruleset/valid/cso.yaml"),
			},
			wantErr: false,
		},
		{
			name: "rule violation",
			fields: fields{
				ContainerReader: cmdutil.FileReader("../../../test/fixtures/bundles/complete.bundle"),
				RuleSetReader:   cmdutil.FileReader("../../../test/fixtures/ruleset/valid/database-secret-validator.yaml"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &LintTask{
				ContainerReader: tt.fields.ContainerReader,
				RuleSetReader:   tt.fields.RuleSetReader,
			}
			if err := tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("LintTask.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
