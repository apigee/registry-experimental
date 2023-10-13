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

package configs

import (
	"context"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/servicebuilder"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"

	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var output string
	var producerProject string
	cmd := &cobra.Command{
		Use:  "create SERVICE",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return nil
			}
			c, err := config.GetConfig()
			if err != nil {
				return nil
			}
			api := &servicebuilder.Api{
				Name:    args[0],
				Version: "1.0.0",
				Operations: []servicebuilder.Operation{
					{
						Id:     "Hello",
						Method: "GET",
						Path:   "/hello",
					},
					{
						Id:     "Goodbye",
						Method: "POST",
						Path:   "/goodbye",
					},
					{
						Id:     "Unknown",
						Method: "GET",
						Path:   "/unknown",
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
					Summary: c.Summary,
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
				ProducerProjectId:  producerProject,
				SystemParameters:   nil,
				Title:              c.Title,
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
			srv, err := servicemanagement.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return nil
			}
			item, err := srv.Services.Configs.
				Create(args[0], service).
				Do()
			if err != nil {
				return nil
			}
			bytes, err := patch.MarshalAndWrap(item, "Operation", item.Name, output)
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
	cmd.Flags().StringVarP(&producerProject, "project", "p", "", "Producer project.")
	return cmd
}
