// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package kv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
)

func TestSecretReader_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/internal/ui/mounts/application/secret/not/found":
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"data":{"type":"kv", "path":"application/", "options":{"version": "2"}}}`)
		case "/v1/application/data/secret/not/found":
			w.WriteHeader(404)
			fmt.Fprintf(w, `{}`)
		default:
			w.WriteHeader(400)
		}
	}))
	defer server.Close()

	// Initialize Vault client
	vaultClient, err := api.NewClient(&api.Config{
		Address:    server.URL,
		Timeout:    time.Second * 1,
		MaxRetries: 1,
		HttpClient: &http.Client{Transport: cleanhttp.DefaultTransport(), Timeout: time.Second * 2},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Build reader
	underTest := SecretGetter(context.Background(), vaultClient)

	_, err = underTest("application/secret/not/found")
	if err != nil && !errors.Is(err, ErrPathNotFound) {
		t.Errorf("SecretReader() error = %v, expected %v", err, ErrPathNotFound)
	}
}

func TestSecretReader_Found(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/internal/ui/mounts/application/secret/found":
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"data":{"type":"kv", "path":"application/", "options":{"version": "2"}}}`)
		case "/v1/application/data/secret/found":
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"data":{"data":{"key":"value"},"metadata":{}}}`)
		default:
			w.WriteHeader(400)
		}
	}))
	defer server.Close()

	// Initialize Vault client
	vaultClient, err := api.NewClient(&api.Config{
		Address:    server.URL,
		Timeout:    time.Second * 1,
		MaxRetries: 1,
		HttpClient: &http.Client{Transport: cleanhttp.DefaultTransport(), Timeout: time.Second * 2},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Build reader
	underTest := SecretGetter(context.Background(), vaultClient)

	res, err := underTest("application/secret/found")
	if err != nil {
		t.Errorf("SecretReader() error = %v, expected nil", err)
	}
	expectedRes := map[string]interface{}{
		"key": "value",
	}
	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("SecretReader() got %v, expected %v", res, expectedRes)
	}
}
