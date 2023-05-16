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

package count

import (
	"context"
	"strconv"
	"testing"

	"github.com/apigee/registry/pkg/connection/grpctest"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/test/seeder"
)

func TestCountDeploymentRevisions(t *testing.T) {
	ctx := context.Background()
	client, _ := grpctest.SetupRegistry(ctx, t, "count-test", []seeder.RegistryResource{
		&rpc.ApiDeployment{Name: "projects/count-test/locations/global/apis/1/deployments/1"},
		&rpc.ApiDeployment{Name: "projects/count-test/locations/global/apis/1/deployments/2"},
	})

	updateDepReq := &rpc.UpdateApiDeploymentRequest{
		ApiDeployment: &rpc.ApiDeployment{
			Name:            "projects/count-test/locations/global/apis/1/deployments/2",
			ApiSpecRevision: "projects/p/apis/a/versions/v/specs/s@12345678",
		},
	}
	_, err := client.UpdateApiDeployment(ctx, updateDepReq)
	if err != nil {
		t.Fatalf("Setup: UpdateApiDeployment(%+v) returned error: %s", updateDepReq, err)
	}

	cmd := Command()
	args := []string{"revisions", "projects/count-test/locations/global/apis/1/deployments/-"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %v returned error: %s", args, err)
	}

	tests := []struct {
		deployment string
		count      int
	}{
		{"projects/count-test/locations/global/apis/1/deployments/1", 1},
		{"projects/count-test/locations/global/apis/1/deployments/2", 2},
	}

	for _, test := range tests {
		t.Run(test.deployment, func(t *testing.T) {
			api, err := client.GetApiDeployment(ctx, &rpc.GetApiDeploymentRequest{
				Name: test.deployment,
			})
			if err != nil {
				t.Fatal("failed GetApiDeployment", err)
			}
			count, err := strconv.Atoi(api.Labels["revisions"])
			if err != nil {
				t.Fatal("failed strconv", err)
			}
			if count != test.count {
				t.Errorf("expected %d, got %d", test.count, count)
			}
		})
	}
}

func TestCountSpecRevisions(t *testing.T) {
	ctx := context.Background()
	client, _ := grpctest.SetupRegistry(ctx, t, "count-test", []seeder.RegistryResource{
		&rpc.ApiSpec{Name: "projects/count-test/locations/global/apis/1/versions/1/specs/1"},
		&rpc.ApiSpec{Name: "projects/count-test/locations/global/apis/1/versions/1/specs/2"},
	})

	updateSpecReq := &rpc.UpdateApiSpecRequest{
		ApiSpec: &rpc.ApiSpec{
			Name:     "projects/count-test/locations/global/apis/1/versions/1/specs/2",
			Contents: []byte(`whatever`),
		},
	}
	_, err := client.UpdateApiSpec(ctx, updateSpecReq)
	if err != nil {
		t.Fatalf("Setup: UpdateApiSpec(%+v) returned error: %s", updateSpecReq, err)
	}

	cmd := Command()
	args := []string{"revisions", "projects/count-test/locations/global/apis/1/versions/1/specs/-"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %v returned error: %s", args, err)
	}

	tests := []struct {
		spec  string
		count int
	}{
		{"projects/count-test/locations/global/apis/1/versions/1/specs/1", 1},
		{"projects/count-test/locations/global/apis/1/versions/1/specs/2", 2},
	}

	for _, test := range tests {
		t.Run(test.spec, func(t *testing.T) {
			api, err := client.GetApiSpec(ctx, &rpc.GetApiSpecRequest{
				Name: test.spec,
			})
			if err != nil {
				t.Fatal("failed GetApiSpec", err)
			}
			count, err := strconv.Atoi(api.Labels["revisions"])
			if err != nil {
				t.Fatal("failed strconv", err)
			}
			if count != test.count {
				t.Errorf("expected %d, got %d", test.count, count)
			}
		})
	}
}
