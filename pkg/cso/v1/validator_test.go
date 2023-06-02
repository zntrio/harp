// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package v1

import "testing"

var tests = []struct {
	in      string
	wantErr bool
}{
	{"", true},
	{"bad", true},
	{"bad/", true},
	{"bad/foo", true},
	// Meta
	{"meta/cso", true},
	{"meta//revision", true},
	{"meta/cso/revision", false},
	// Infra
	{"infra", true},
	{"infra/", true},
	{"infra/foo", true},
	{"infra/foo//", true},
	{"infra/aws//invalid-region/iam", true},
	{"infra/aws/security/invalid-region/iam", true},
	{"infra/aws/security/us-east-15/rds/postgres/admin_creds", true},
	{"infra/gcp/security/us-east-1/db/postgres/admin_creds", true},
	{"infra/gcp/security/us-east15/db/postgres/admin_creds", true},
	{"infra/aws/security/global/iam", false},
	{"infra/gcp/security/global/iam", false},
	{"infra/azure/security/global/iam", false},
	{"infra/aws/security/us-east-1/rds/postgres/admin_creds", false},
	{"infra/gcp/security/us-east1/db/postgres/admin_creds", false},
	{"infra/ibm/security/au-syd/db/postgres/admin_creds", false},
	{"infra/azure-us-gov/security/usgovvirginia/db/postgres/admin_creds", false},
	{"infra/local/security/global/dns/registrar/admin_creds", false},
	{"infra/unsupported/security/global/dns/registrar/admin_creds", true},
	// Platform
	{"platform", true},
	{"platform/", true},
	{"platform/foo", true},
	{"platform/foo//", true},
	{"platform/production", true},
	{"platform/staging", true},
	{"platform/qa", true},
	{"platform/dev", true},
	{"platform/dev/foo", true},
	{"platform/production/foo/eu-central-1", true},
	{"platform/production/foo/eu-central-1/database", true},
	{"platform/production/foo//", true},
	{"platform/production////test", true},
	{"platform/invalid/foo/eu-central-1/db//admin_account", true},
	{"platform/production/foo/invalid-region/db/admin_account", true},
	{"platform/production/foo/invalid-region/db/admin_account", true},
	{"platform/production/foo/eu-central-1//admin_account", true},
	{"platform/production/foo/eu-central-1/db/admin_account", false},
	// Product
	{"product", true},
	{"product/", true},
	{"product/foo", true},
	{"product/foo//", true},
	{"product/foo", true},
	{"product/foo/v1.0.0", true},
	{"product//v1.0.0/foo", true},
	{"product/foo/abc/foo", true},
	{"product/foo/v1.0.0/foo", false},
	{"product/foo/v1.0.0/foo/bar", false},
	{"product/foo/1.0.0/foo", false},
	{"product/foo/1.0.0/foo/bar", false},
	// Application
	{"app", true},
	{"app/production/name/foo/v1.0.0//foo", true},
	{"app/production/name/foo/v1.0.0/component/foo", false},
	{"app/production/name/foo/v1.0.0/component/foo/bar", false},
	{"app/essp/name/foo/v1.0.0/component/foo/bar", true},
	// Artifact
	{"artifact", true},
	{"artifact/docker", true},
	{"artifact/docker/sha256:fab2dded59dd0c2894dd9dbae71418f565be5bd0d8fd82365c16aec41c7e367f/attestations/snyk_report", false},
}

func Test_Validate(t *testing.T) {
	for _, tt := range tests {
		err := Validate(tt.in)
		if tt.wantErr != (err != nil) {
			t.Errorf("Validate(%q) = %v, want %v", tt.in, err, tt.wantErr)
		}
	}
}
