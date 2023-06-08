// Copyright 2023 Google LLC.
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

func TestCountDeployments(t *testing.T) {
	ctx := context.Background()
	client, _ := grpctest.SetupRegistry(ctx, t, "count-test", []seeder.RegistryResource{
		&rpc.ApiDeployment{
			Name: "projects/count-test/locations/global/apis/1/deployments/1",
		},
		&rpc.ApiDeployment{
			Name: "projects/count-test/locations/global/apis/2/deployments/1",
		},
		&rpc.ApiDeployment{
			Name: "projects/count-test/locations/global/apis/2/deployments/2",
		},
	})

	cmd := Command()
	args := []string{"deployments", "projects/count-test/locations/global/apis/-"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() with args %v returned error: %s", args, err)
	}

	tests := []struct {
		api   string
		count int
	}{
		{"projects/count-test/locations/global/apis/1", 1},
		{"projects/count-test/locations/global/apis/2", 2},
	}

	for _, test := range tests {
		t.Run(test.api, func(t *testing.T) {
			api, err := client.GetApi(ctx, &rpc.GetApiRequest{
				Name: test.api,
			})
			if err != nil {
				t.Fatal("failed GetApi", err)
			}
			count, err := strconv.Atoi(api.Labels["deployments"])
			if err != nil {
				t.Fatal("failed strconv", err)
			}
			if count != test.count {
				t.Errorf("expected %d, got %d", test.count, count)
			}
		})
	}
}
