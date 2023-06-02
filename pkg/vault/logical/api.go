// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package logical

import "github.com/hashicorp/vault/api"

//go:generate mockgen -destination logical.mock.go -package logical zntr.io/harp/v2/pkg/vault/logical Logical

// Logical backend interface.
type Logical interface {
	Read(path string) (*api.Secret, error)
	ReadWithData(path string, data map[string][]string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
	WriteBytes(path string, data []byte) (*api.Secret, error)
	List(path string) (*api.Secret, error)
	Unwrap(token string) (*api.Secret, error)
	Delete(path string) (*api.Secret, error)
	DeleteWithData(path string, data map[string][]string) (*api.Secret, error)
}
