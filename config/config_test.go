// Copyright 2019 The Morning Consult, LLC or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//         https://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	vaultconfig "github.com/hashicorp/vault/command/agent/config"
)

func TestLoadConfig(t *testing.T) {

	cases := []struct {
		name         string
		file         string
		err          string
		expectConfig *vaultconfig.Config
	}{
		{
			"file-doesnt-exist",
			"testdata/nonexistent.hcl",
			"stat testdata/nonexistent.hcl: no such file or directory",
			nil,
		},
		{
			"provided-directory",
			"testdata",
			"location is a directory, not a file",
			nil,
		},
		{
			"empty-file",
			"testdata/empty-file.hcl",
			"no 'auto_auth' block found in configuration file",
			nil,
		},
		{
			"no-method",
			"testdata/no-method.hcl",
			"error parsing 'auto_auth': error parsing 'method': one and only one \"method\" block is required",
			nil,
		},
		{
			"no-sinks",
			"testdata/no-sinks.hcl",
			"",
			&vaultconfig.Config{
				AutoAuth: &vaultconfig.AutoAuth{
					Method: &vaultconfig.Method{
						Type:      "approle",
						MountPath: "auth/approle",
						Config: map[string]interface{}{
							"role_id_file_path":   "/tmp/role-id",
							"secret":              "secret/docker/creds",
							"secret_id_file_path": "/tmp/secret-id",
						},
					},
				},
			},
		},
		{
			"no-mount-path",
			"testdata/no-mount-path.hcl",
			"",
			&vaultconfig.Config{
				AutoAuth: &vaultconfig.AutoAuth{
					Method: &vaultconfig.Method{
						Type:      "aws",
						MountPath: "auth/aws",
						Config: map[string]interface{}{
							"role":   "dev-role-iam",
							"secret": "secret/docker/creds",
							"type":   "iam",
						},
					},
					Sinks: []*vaultconfig.Sink{
						{
							Type: "file",
							Config: map[string]interface{}{
								"path": "/tmp/foo",
							},
						},
					},
				},
			},
		},
		{
			"valid",
			"testdata/valid.hcl",
			"",
			&vaultconfig.Config{
				AutoAuth: &vaultconfig.AutoAuth{
					Method: &vaultconfig.Method{
						Type:      "approle",
						MountPath: "auth/approle",
						Config: map[string]interface{}{
							"role_id_file_path":   "/tmp/role-id",
							"secret":              "secret/docker/creds",
							"secret_id_file_path": "/tmp/secret-id",
						},
					},
					Sinks: []*vaultconfig.Sink{
						{
							Type: "file",
							Config: map[string]interface{}{
								"path": "/tmp/foo",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotConfig, err := LoadConfig(tc.file)
			if tc.err != "" {
				if err == nil {
					t.Fatal("expected an error but didn't receive one")
				}
				if err.Error() != tc.err {
					t.Fatalf("Results differ:\n%v", cmp.Diff(err.Error(), tc.err))
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			comparer := cmp.Comparer(func(c1 *vaultconfig.Config, c2 *vaultconfig.Config) bool {
				if (c1 == nil || c2 == nil) && !((c1 == nil) && (c2 == nil)) {
					return false
				}
				return cmp.Equal(c1.AutoAuth, c2.AutoAuth) && cmp.Equal(c1.Vault, c2.Vault)
			})
			if !cmp.Equal(tc.expectConfig, gotConfig, comparer) {
				t.Errorf("Configurations differ:\n%v", cmp.Diff(tc.expectConfig, gotConfig))
			}
		})
	}
}