// Copyright 2023 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bleve

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/pkg/connection/grpctest"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry"
	"github.com/apigee/registry/server/registry/test/seeder"
)

func TestMain(m *testing.M) {
	grpctest.TestMain(m, registry.Config{})
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	blevePath := filepath.Join(t.TempDir(), "registry.bleve")

	bookstore_protos, err := compress.ZipArchiveOfPath("testdata/examples", "testdata/", true)
	if err != nil {
		t.Fatalf("Failed to initialize search test: %s", err)
	}
	discovery_discovery, err := os.ReadFile("testdata/discovery-v1.json")
	if err != nil {
		t.Fatalf("Failed to initialize search test: %s", err)
	}
	discovery_discovery, err = compress.GZippedBytes(discovery_discovery)
	if err != nil {
		t.Fatalf("Failed to initialize search test: %s", err)
	}
	petstore_openapi, err := os.ReadFile("testdata/petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to initialize search test: %s", err)
	}
	grpctest.SetupRegistry(ctx, t, "search-test",
		[]seeder.RegistryResource{
			&rpc.ApiSpec{
				Name:     "projects/search-test/locations/global/apis/bookstore/versions/v1/specs/protos",
				MimeType: "application/x.protos+zip",
				Filename: "protos.zip",
				Contents: bookstore_protos.Bytes(),
			},
			&rpc.ApiSpec{
				Name:     "projects/search-test/locations/global/apis/discovery/versions/v1/specs/discovery",
				MimeType: "application/x.discovery+gzip",
				Filename: "discovery.json",
				Contents: discovery_discovery,
			},
			&rpc.ApiSpec{
				Name:     "projects/search-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi",
				MimeType: "application/x.openapi;version=3.0.0",
				Filename: "openapi.yaml",
				Contents: petstore_openapi,
			},
			&rpc.ApiSpec{
				Name:     "projects/search-test/locations/global/apis/hello/versions/v1/specs/plain",
				MimeType: "text/plain",
				Filename: "hello.text",
				Contents: []byte("Hello, this is an http API."),
			},
		})

	cmd := Command()
	cmd.SetArgs([]string{"index", "projects/search-test/locations/global/apis/-/versions/-/specs/-", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}

	cmd = Command()
	cmd.SetArgs([]string{"search", "books", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}

	cmd = Command()
	cmd.SetArgs([]string{"search", "discovery", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}

	cmd = Command()
	cmd.SetArgs([]string{"search", "pets", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}

	cmd = Command()
	cmd.SetArgs([]string{"search", "http", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}
}
