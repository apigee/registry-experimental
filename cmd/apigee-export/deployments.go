// Copyright 2022 Google LLC. All Rights Reserved.
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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apigee/registry/log"
	"github.com/apigee/registry/pkg/models"
	"github.com/spf13/cobra"
	"google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v2"
)

var exportDeploymentsCommand = &cobra.Command{
	Use:   "deployments ORGANIZATION [DIRECTORY]",
	Short: "Exports Apigee deployments to YAML files compatible with API Registry",
	Args:  cobra.RangeArgs(1, 2),
	Run:   exportDeployments,
}

func exportDeployments(cmd *cobra.Command, args []string) {
	var (
		ctx = cmd.Context()
		org = args[0]
	)

	env, err := newEnvMap(ctx, org)
	if err != nil {
		log.Fatalf(ctx, "Failed to get hostnames for environments in %s: %s", org, err)
	}

	deps, err := deployments(ctx, org)
	if err != nil {
		log.Fatalf(ctx, "Failed to list deployments for %s: %s", org, err)
	}

	for _, dep := range deps {
		api := &models.Api{
			Header: models.Header{
				ApiVersion: "apigeeregistry/v1",
				Kind:       "API",
				Metadata: models.Metadata{
					Name: clean(dep.ApiProxy),
					Annotations: map[string]string{
						"apigee-organization": org,
						"apigeex-proxy":       fmt.Sprintf("%s/apis/%s", org, dep.ApiProxy),
					},
				},
			},
			Data: models.ApiData{
				DisplayName:    dep.ApiProxy,
				ApiDeployments: make([]*models.ApiDeployment, 0, len(deps)),
			},
		}

		hostnames, ok := env.Hostnames(dep.Environment)
		if !ok {
			log.Warnf(ctx, "Failed to find hostnames for environment %s", dep.Environment)
			continue
		}

		for _, hostname := range hostnames {
			envgroup, ok := env.Envgroup(hostname)
			if !ok {
				log.Warnf(ctx, "Failed to determine envgroup for hostname %q", hostname)
				continue
			}

			api.Data.ApiDeployments = append(api.Data.ApiDeployments, &models.ApiDeployment{
				Header: models.Header{
					ApiVersion: "apigeeregistry/v1",
					Kind:       "Deployment",
					Metadata: models.Metadata{
						Name: clean(hostname),
						Annotations: map[string]string{
							"apigee-organization":   org,
							"apigee-proxy-revision": fmt.Sprintf("%s/apis/%s/revisions/%s", org, dep.ApiProxy, dep.Revision),
							"apigeex-environment":   fmt.Sprintf("%s/environments/%s", org, dep.Environment),
							"apigee-envgroup":       envgroup,
						},
					},
				},
				Data: models.ApiDeploymentData{
					DisplayName: dep.Environment,
					EndpointURI: hostname,
				},
			})
		}

		out, err := yaml.Marshal(api)
		if err != nil {
			log.Errorf(ctx, "Failed to marshal YAML for model: %s", err)
			continue
		}

		if verbose {
			fmt.Println(string(out))
		}

		// Only write the files if a directory is specified.
		if len(args) < 2 {
			continue
		}

		filename := filepath.Join(args[1], api.Metadata.Name+".yaml")
		if err := os.WriteFile(filename, out, 0644); err != nil {
			log.Errorf(ctx, "Failed to write YAML for API: %s", err)
		}
	}
}

func deployments(ctx context.Context, org string) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Deployments.List(org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Deployments, nil
}
