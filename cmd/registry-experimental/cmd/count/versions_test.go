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

func TestCountVersions(t *testing.T) {
	ctx := context.Background()
	client, _ := grpctest.SetupRegistry(ctx, t, "count-test", []seeder.RegistryResource{
		&rpc.ApiVersion{
			Name: "projects/count-test/locations/global/apis/1/versions/1",
		},
		&rpc.ApiVersion{
			Name: "projects/count-test/locations/global/apis/2/versions/1",
		},
		&rpc.ApiVersion{
			Name: "projects/count-test/locations/global/apis/2/versions/2",
		},
	})

	cmd := Command()
	args := []string{"versions", "projects/count-test/locations/global/apis/-"}
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
			count, err := strconv.Atoi(api.Labels["versions"])
			if err != nil {
				t.Fatal("failed strconv", err)
			}
			if count != test.count {
				t.Errorf("expected %d, got %d", test.count, count)
			}
		})
	}
}
