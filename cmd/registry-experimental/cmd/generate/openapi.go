// Copyright 2021 Google LLC. All Rights Reserved.
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

package generate

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"github.com/spf13/cobra"
)

func openapiCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Generate an OpenAPI spec for a protocol buffer API specification",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get filter from flags")
			}

			client, err := connection.NewClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			taskQueue, wait := core.WorkerPool(ctx, 16)
			defer wait()

			// Generate tasks.
			name := args[0]
			if spec, err := names.ParseSpec(name); err == nil {
				// Iterate through a collection of specs and evaluate each.
				err = core.ListSpecs(ctx, client, spec, filter, func(spec *rpc.ApiSpec) {
					taskQueue <- &generateOpenAPITask{
						client:   client,
						specName: spec.Name,
					}
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list specs")
				}
			}
		},
	}
}

type generateOpenAPITask struct {
	client   connection.Client
	specName string
}

func (task *generateOpenAPITask) String() string {
	return fmt.Sprintf("generate openapi for %s", task.specName)
}

func (task *generateOpenAPITask) Run(ctx context.Context) error {
	log.FromContext(ctx).Info(task.String())
	return nil
}
