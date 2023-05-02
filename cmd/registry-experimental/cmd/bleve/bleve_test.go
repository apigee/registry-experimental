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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/pkg/connection/grpctest"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry"
	"github.com/apigee/registry/server/registry/test/seeder"
)

func TestMain(m *testing.M) {
	grpctest.TestMain(m, registry.Config{})
}

func buildTestRegistry(ctx context.Context, t *testing.T) {
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
}

func TestSearch(t *testing.T) {
	var err error
	ctx := context.Background()
	blevePath := filepath.Join(t.TempDir(), "registry.bleve")

	buildTestRegistry(ctx, t)

	cmd := Command()
	cmd.SetArgs([]string{"index", "projects/search-test/locations/global/apis/-/versions/-/specs/-", "--bleve", blevePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
	}

	type searchResult struct {
		TotalHits int `json:"total_hits"`
	}

	t.Run("search-books", func(t *testing.T) {
		cmd := Command()
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"search", "books", "--bleve", blevePath, "--output", "json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		var result searchResult
		err = json.Unmarshal(buf.Bytes(), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal search results: %s", err)
		}
		if result.TotalHits != 1 {
			t.Errorf("Expected 1 hit, got %d", result.TotalHits)
		}
	})

	t.Run("search-discovery", func(t *testing.T) {
		cmd := Command()
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"search", "discovery", "--bleve", blevePath, "--output", "json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		var result searchResult
		err = json.Unmarshal(buf.Bytes(), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal search results: %s", err)
		}
		if result.TotalHits != 1 {
			t.Errorf("Expected 1 hit, got %d", result.TotalHits)
		}
	})

	t.Run("search-pets", func(t *testing.T) {
		cmd := Command()
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"search", "pets", "--bleve", blevePath, "--output", "json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		var result searchResult
		err = json.Unmarshal(buf.Bytes(), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal search results: %s", err)
		}
		if result.TotalHits != 1 {
			t.Errorf("Expected 1 hit, got %d", result.TotalHits)
		}
	})

	t.Run("search-http", func(t *testing.T) {
		cmd := Command()
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"search", "http", "--bleve", blevePath, "--output", "json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		var result searchResult
		err = json.Unmarshal(buf.Bytes(), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal search results: %s", err)
		}
		if result.TotalHits != 4 {
			t.Errorf("Expected 4 hits, got %d", result.TotalHits)
		}
	})

	t.Run("search-nohits", func(t *testing.T) {
		cmd := Command()
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"search", "nohits", "--bleve", blevePath, "--output", "json"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		var result searchResult
		err = json.Unmarshal(buf.Bytes(), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal search results: %s", err)
		}
		if result.TotalHits != 0 {
			t.Errorf("Expected 0 hits, got %d", result.TotalHits)
		}
	})
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	blevePath := filepath.Join(t.TempDir(), "registry.bleve")
	port := "8891"

	buildTestRegistry(ctx, t)

	// Start the server.
	go func() {
		cmd := Command()
		cmd.SetArgs([]string{"serve", "--bleve", blevePath, "--port", port})
		if err := cmd.Execute(); err != nil {
			log.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
	}()

	// Wait for the server to start.
	_, err := net.DialTimeout("tcp", "localhost:"+port, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to test server: %s", err)
	}

	// Call the indexing API.
	t.Run("server-indexing", func(t *testing.T) {
		postBody, err := json.Marshal(struct {
			Pattern string `json:"pattern"`
			Filter  string `json:"string"`
		}{
			Pattern: "projects/search-test/locations/global/apis/-/versions/-/specs/-",
			Filter:  "mime_type.contains('openapi')",
		})
		if err != nil {
			t.Fatalf("failed to create request %s", err)
		}
		request, err := http.NewRequest("POST", "http://localhost:"+port+"/index", bytes.NewBuffer(postBody))
		if err != nil {
			t.Fatalf("failed to create request %s", err)
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("failed to call index API %s", err)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			t.Fatalf("unexpected code from index API %d", response.StatusCode)
		}
		fmt.Println("response Headers:", response.Header)
	})

	// Call the search API.
	t.Run("server-search", func(t *testing.T) {
		request, err := http.NewRequest("GET", "http://localhost:"+port+"/search?q=pet", nil)
		if err != nil {
			t.Fatalf("failed to create request %s", err)
		}
		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			t.Fatalf("failed to call search API %s", err)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			t.Fatalf("unexpected code from search API %d", response.StatusCode)
		}
		fmt.Println("response Headers:", response.Header)
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("failed to read search API response %s", err)
		}
		var searchResponse struct {
			TotalHits int `json:"total_hits"`
		}
		if err = json.Unmarshal(body, &searchResponse); err != nil {
			t.Fatalf("failed to read unmarshal API response %s", err)
		}
		if searchResponse.TotalHits != 1 {
			t.Fatalf("failed to get expected number of hits (1), got %d", searchResponse.TotalHits)
		}
	})
}
