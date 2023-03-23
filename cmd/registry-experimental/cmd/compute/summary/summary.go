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
	"fmt"

	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary PATTERN",
		Short: "Compute summary metrics of resources",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				return err
			}
			pattern := c.FQName(args[0])
			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				return err
			}
			adminClient, err := connection.NewAdminClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			registryClient, err := connection.NewRegistryClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			v := &summaryVisitor{
				registryClient: registryClient,
			}
			return visitor.Visit(ctx, v, visitor.VisitorOptions{
				RegistryClient:  registryClient,
				AdminClient:     adminClient,
				Pattern:         pattern,
				Filter:          filter,
				ImplicitProject: &rpc.Project{Name: "projects/implicit"},
			})
		},
	}
	cmd.Flags().String("filter", "", "Filter selected resources")
	cmd.Flags().Int("jobs", 10, "Number of actions to perform concurrently")
	return cmd
}

type summaryVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
}

type Summary struct {
	ApiCount        int            `yaml:"apis,omitempty"`
	VersionCount    int            `yaml:"versions,omitempty"`
	SpecCount       int            `yaml:"specs,omitempty"`
	DeploymentCount int            `yaml:"deployments,omitempty"`
	MimeTypes       map[string]int `yaml:"mimetypes"`
}

func (v *summaryVisitor) ProjectHandler() visitor.ProjectHandler {
	return func(ctx context.Context, message *rpc.Project) error {
		fmt.Printf("%s\n", message.Name)
		projectName, err := names.ParseProject(message.Name)
		if err != nil {
			return err
		}
		apiCount := 0
		if err := visitor.ListAPIs(ctx, v.registryClient,
			projectName.Api("-"), "",
			func(ctx context.Context, api *rpc.Api) error {
				apiCount++
				return nil
			}); err != nil {
			return err
		}
		versionCount := 0
		if err := visitor.ListVersions(ctx, v.registryClient,
			projectName.Api("-").Version("-"), "",
			func(ctx context.Context, message *rpc.ApiVersion) error {
				versionCount++
				return nil
			}); err != nil {
			return err
		}
		specCount := 0
		mimeTypes := make(map[string]int)
		if err := visitor.ListSpecs(ctx, v.registryClient,
			projectName.Api("-").Version("-").Spec("-"), "", false,
			func(ctx context.Context, message *rpc.ApiSpec) error {
				specCount++
				mimeTypes[message.MimeType]++
				return nil
			}); err != nil {
			return err
		}
		deploymentCount := 0
		if err := visitor.ListDeployments(ctx, v.registryClient,
			projectName.Api("-").Deployment("-"), "",
			func(ctx context.Context, message *rpc.ApiDeployment) error {
				deploymentCount++
				return nil
			}); err != nil {
			return err
		}
		summary := &Summary{
			ApiCount:        apiCount,
			VersionCount:    versionCount,
			SpecCount:       specCount,
			DeploymentCount: deploymentCount,
			MimeTypes:       mimeTypes,
		}
		bytes, err := yaml.Marshal(summary)
		if err != nil {
			return err
		}
		artifact := &rpc.Artifact{
			Name:     projectName.Artifact("summary").String(),
			MimeType: "application/yaml;type=Summary",
			Contents: bytes,
		}
		return visitor.SetArtifact(ctx, v.registryClient, artifact)
	}
}

func (v *summaryVisitor) ApiHandler() visitor.ApiHandler {
	return func(ctx context.Context, message *rpc.Api) error {
		fmt.Printf("%s\n", message.Name)
		apiName, err := names.ParseApi(message.Name)
		if err != nil {
			return err
		}
		versionCount := 0
		if err := visitor.ListVersions(ctx, v.registryClient,
			apiName.Version("-"), "",
			func(ctx context.Context, message *rpc.ApiVersion) error {
				versionCount++
				return nil
			}); err != nil {
			return err
		}
		specCount := 0
		mimeTypes := make(map[string]int)
		if err := visitor.ListSpecs(ctx, v.registryClient,
			apiName.Version("-").Spec("-"), "", false,
			func(ctx context.Context, message *rpc.ApiSpec) error {
				specCount++
				mimeTypes[message.MimeType]++
				return nil
			}); err != nil {
			return err
		}
		deploymentCount := 0
		if err := visitor.ListDeployments(ctx, v.registryClient,
			apiName.Deployment("-"), "",
			func(ctx context.Context, message *rpc.ApiDeployment) error {
				deploymentCount++
				return nil
			}); err != nil {
			return err
		}
		summary := &Summary{
			ApiCount:        1,
			VersionCount:    versionCount,
			SpecCount:       specCount,
			DeploymentCount: deploymentCount,
			MimeTypes:       mimeTypes,
		}
		bytes, err := yaml.Marshal(summary)
		if err != nil {
			return err
		}
		artifact := &rpc.Artifact{
			Name:     apiName.Artifact("summary").String(),
			MimeType: "application/yaml;type=Summary",
			Contents: bytes,
		}
		return visitor.SetArtifact(ctx, v.registryClient, artifact)
	}
}
