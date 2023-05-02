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

package summary

import (
	"context"
	"testing"

	"github.com/apigee/registry/pkg/connection/grpctest"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry"
	"github.com/apigee/registry/server/registry/test/seeder"
	"gopkg.in/yaml.v2"
)

func TestMain(m *testing.M) {
	grpctest.TestMain(m, registry.Config{})
}

func TestSummary(t *testing.T) {
	ctx := context.Background()
	registryClient, _ := grpctest.SetupRegistry(ctx, t, "summary-test",
		[]seeder.RegistryResource{
			&rpc.ApiSpec{
				Name:     "projects/summary-test/locations/global/apis/a1/versions/v1/specs/s1",
				MimeType: "text/plain",
			},
			&rpc.ApiSpec{
				Name:     "projects/summary-test/locations/global/apis/a2/versions/v1/specs/s1",
				MimeType: "text/plain",
			},
			&rpc.ApiSpec{
				Name:     "projects/summary-test/locations/global/apis/a2/versions/v1/specs/s2",
				MimeType: "application/json",
			},
			&rpc.Artifact{
				Name: "projects/summary-test/locations/global/apis/a1/versions/v1/specs/s1/artifacts/x1",
			},
			&rpc.Artifact{
				Name: "projects/summary-test/locations/global/apis/a2/versions/v1/specs/s1/artifacts/x1",
			},
			&rpc.Artifact{
				Name: "projects/summary-test/locations/global/apis/a2/versions/v1/specs/s2/artifacts/x1",
			},
			&rpc.Artifact{
				Name: "projects/summary-test/locations/global/apis/a2/deployments/d1/artifacts/x1",
			},
		})

	t.Run("project-summary", func(t *testing.T) {
		cmd := Command()
		cmd.SetArgs([]string{"projects/summary-test"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		contents, err := registryClient.GetArtifactContents(ctx,
			&rpc.GetArtifactContentsRequest{Name: "projects/summary-test/locations/global/artifacts/summary"})
		if err != nil {
			t.Fatalf("Artifact could not be read %s", err)
		}
		var s Summary
		if err := yaml.Unmarshal(contents.Data, &s); err != nil {
			t.Fatalf("Artifact could not be parsed %s", err)
		}
		if s.ApiCount != 2 {
			t.Errorf("Incorrect API count, expected 2, got %d", s.ApiCount)
		}
		if s.VersionCount != 2 {
			t.Errorf("Incorrect version count, expected 2, got %d", s.VersionCount)
		}
		if s.SpecCount != 3 {
			t.Errorf("Incorrect spec count, expected 3, got %d", s.SpecCount)
		}
		if s.DeploymentCount != 1 {
			t.Errorf("Incorrect deployment count, expected 1, got %d", s.DeploymentCount)
		}
		if len(s.MimeTypes) != 2 {
			t.Errorf("Incorrect mime_type count, expected 2, got %d", len(s.MimeTypes))
		}
	})

	t.Run("api-summary", func(t *testing.T) {
		cmd := Command()
		cmd.SetArgs([]string{"projects/summary-test/locations/global/apis/a2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", cmd.Args, err)
		}
		contents, err := registryClient.GetArtifactContents(ctx,
			&rpc.GetArtifactContentsRequest{Name: "projects/summary-test/locations/global/apis/a2/artifacts/summary"})
		if err != nil {
			t.Fatalf("Artifact could not be read %s", err)
		}
		var s Summary
		if err := yaml.Unmarshal(contents.Data, &s); err != nil {
			t.Fatalf("Artifact could not be parsed %s", err)
		}
		if s.ApiCount != 1 {
			t.Errorf("Incorrect API count, expected 1, got %d", s.ApiCount)
		}
		if s.VersionCount != 1 {
			t.Errorf("Incorrect version count, expected 1, got %d", s.VersionCount)
		}
		if s.SpecCount != 2 {
			t.Errorf("Incorrect spec count, expected 2, got %d", s.SpecCount)
		}
		if s.DeploymentCount != 1 {
			t.Errorf("Incorrect deployment count, expected 1, got %d", s.DeploymentCount)
		}
		if len(s.MimeTypes) != 2 {
			t.Errorf("Incorrect mime_type count, expected 2, got %d", len(s.MimeTypes))
		}
	})
}
