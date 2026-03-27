// Copyright 2023 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/servicebuilder"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"

	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	var apply bool
	var output string
	var project string
	cmd := &cobra.Command{
		Use:  "config SERVICE",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			api := &servicebuilder.Api{
				Name:    args[0],
				Version: "1.0.0",
				Operations: []servicebuilder.Operation{
					{
						Id:     "GetPets",
						Method: "GET",
						Path:   "/v1/pets",
					},
					{
						Id:     "GetPetById",
						Method: "GET",
						Path:   "/v1/pets/:id",
					},
					{
						Id:     "CreatePet",
						Method: "POST",
						Path:   "/v1/pets",
					},
				},
			}
			apis := []*servicebuilder.Api{api}
			service := &servicemanagement.Service{
				Apis:           servicebuilder.Apis(apis),
				Authentication: nil,
				ConfigVersion:  3,
				Control:        servicebuilder.Control(),
				Documentation: &servicemanagement.Documentation{
					Summary: "Petstore API",
				},
				Endpoints: []*servicemanagement.Endpoint{
					{
						Name: args[0],
					},
				},
				Enums:              servicebuilder.Enums(),
				Http:               servicebuilder.Http(apis),
				Logging:            servicebuilder.Logging(),
				Logs:               servicebuilder.Logs(),
				Name:               args[0],
				ProducerProjectId:  project,
				SystemParameters:   nil,
				Title:              "Petstore",
				Metrics:            servicebuilder.Metrics(),
				MonitoredResources: servicebuilder.MonitoredResources(),
				Monitoring:         servicebuilder.Monitoring(),
				// https://cloud.google.com/endpoints/docs/grpc/quotas-configure
				Quota: &servicemanagement.Quota{
					Limits: []*servicemanagement.QuotaLimit{
						{
							Name:   "calls",
							Metric: "calls",
							Unit:   "1/min/{project}",
							Values: map[string]string{
								"STANDARD": "60",
							},
						},
					},
					MetricRules: []*servicemanagement.MetricRule{
						{
							MetricCosts: map[string]string{
								"calls": "1",
							},
							Selector: "*",
						},
					},
				},
				Types: servicebuilder.Types(),
				Usage: servicebuilder.Usage(apis),
			}
			if !apply {
				serviceJSON, err := json.MarshalIndent(service, "", "  ")
				if err != nil {
					return err
				}
				fmt.Printf("%s", string(serviceJSON))
				return nil
			}
			srv, err := servicemanagement.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			item, err := srv.Services.Configs.
				Create(args[0], service).
				Do()
			if err != nil {
				return err
			}
			bytes, err := patch.MarshalAndWrap(item, "ServiceConfig", item.Name, output)
			if err != nil {
				return err
			}
			if _, err := cmd.OutOrStdout().Write(bytes); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of: (yaml, json).")
	cmd.Flags().StringVarP(&project, "project", "p", "", "Project.")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply config.")
	return cmd
}
