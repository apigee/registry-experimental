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

package index

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/apigee/registry-experimental/pkg/yamlquery"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/api/googleapi"
	"gopkg.in/yaml.v3"
)

func operationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operations PATTERN",
		Short: "Build a BigQuery index of API operations",
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
			v := &operationsVisitor{
				registryClient: registryClient,
			}
			client, err := bigquery.NewClient(ctx, "TODO")
			if err != nil {
				return err
			}
			dataset := client.Dataset("registry")
			if err := dataset.Create(ctx, nil); err != nil {
				switch v := err.(type) {
				case *googleapi.Error:
					if v.Code != 409 { // already exists
						return err
					}
				default:
					return err
				}
			}
			table := dataset.Table("operations")

			schema, err := bigquery.InferSchema(operation{})
			if err != nil {
				return err
			}
			if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
				switch v := err.(type) {
				case *googleapi.Error:
					if v.Code != 409 { // already exists
						return err
					}
				default:
					return err
				}
			}
			err = visitor.Visit(ctx, v, visitor.VisitorOptions{
				RegistryClient:  registryClient,
				AdminClient:     adminClient,
				Pattern:         pattern,
				Filter:          filter,
				ImplicitProject: &rpc.Project{Name: "projects/implicit"},
			})
			if err != nil {
				return err
			}
			u := table.Inserter()
			if err := u.Put(ctx, v.operations); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("filter", "", "Filter selected resources")
	return cmd
}

type operation struct {
	Path   string
	Method string
	Spec   string
}

type operationsVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
	operations     []*operation
}

func (v *operationsVisitor) SpecHandler() visitor.SpecHandler {
	return func(ctx context.Context, message *rpc.ApiSpec) error {
		fmt.Printf("%s\n", message.Name)
		specName, err := names.ParseSpec(message.Name)
		if err != nil {
			return err
		}
		return visitor.GetSpec(ctx, v.registryClient, specName, true,
			func(ctx context.Context, spec *rpc.ApiSpec) error {
				if mime.IsOpenAPIv2(spec.MimeType) || mime.IsOpenAPIv3(spec.MimeType) {
					err := v.getOpenAPIOperations(spec.Name, spec.Contents)
					if err != nil {
						return err
					}
				}
				return nil
			})
	}
}

func (v *operationsVisitor) getOpenAPIOperations(specName string, b []byte) error {
	var doc yaml.Node
	err := yaml.Unmarshal(b, &doc)
	if err != nil {
		return err
	}
	paths := yamlquery.QueryNode(&doc, "paths")
	if paths != nil {
		for i := 0; i < len(paths.Content); i += 2 {
			path := paths.Content[i].Value
			methods := paths.Content[i+1]
			for j := 0; j < len(methods.Content); j += 2 {
				method := strings.ToUpper(methods.Content[j].Value)
				if strings.HasPrefix(method, "X-") {
					continue // skip OpenAPI extensions
				}
				fmt.Printf("%s %s\n", method, path)
				v.operations = append(v.operations,
					&operation{
						Method: method,
						Path:   path,
						Spec:   specName,
					})
			}
		}
	}
	return nil
}
