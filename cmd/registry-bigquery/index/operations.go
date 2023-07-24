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
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/apigee/registry-experimental/pkg/yamlquery"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func operationsCommand() *cobra.Command {
	var filter string
	var project string
	var dataset string
	var batchSize int
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
			adminClient, err := connection.NewAdminClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			registryClient, err := connection.NewRegistryClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			if project == "" {
				project = c.Project
			}
			client, err := bigquery.NewClient(ctx, project)
			if err != nil {
				return err
			}
			ds, err := getOrCreateDataset(ctx, client, dataset)
			if err != nil {
				return err
			}
			table, err := getOrCreateTable(ctx, ds, "operations", operation{})
			if err != nil {
				return err
			}
			v := &operationsVisitor{
				registryClient: registryClient,
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
			log.Printf("uploading %d operations", len(v.operations))
			for start := 0; start < len(v.operations); start += batchSize {
				log.Printf("%d", start)
				end := min(start+batchSize, len(v.operations))
				if err := u.Put(ctx, v.operations[start:end]); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "Filter selected resources")
	cmd.Flags().StringVar(&project, "project", "", "Project to use for BigQuery upload (defaults to registry project)")
	cmd.Flags().StringVar(&dataset, "dataset", "registry", "BigQuery dataset name")
	cmd.Flags().IntVar(&batchSize, "batch-size", 10000, "Batch size to use when uploading records to BigQuery")
	return cmd
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type operation struct {
	Path      string
	Method    string
	Api       string
	Version   string
	Spec      string
	Timestamp time.Time
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
					err := v.getOpenAPIOperations(specName, spec.Contents)
					if err != nil {
						return err
					}
				}
				return nil
			})
	}
}

func (v *operationsVisitor) getOpenAPIOperations(specName names.Spec, b []byte) error {
	var doc yaml.Node
	err := yaml.Unmarshal(b, &doc)
	if err != nil {
		return err
	}
	paths := yamlquery.QueryNode(&doc, "paths")
	if paths != nil {
		for i := 0; i < len(paths.Content); i += 2 {
			path := paths.Content[i].Value
			fields := paths.Content[i+1]
			for j := 0; j < len(fields.Content); j += 2 {
				fieldName := fields.Content[j].Value
				// Skip any fields (summary, description, etc) that aren't methods.
				if fieldName != "get" &&
					fieldName != "put" &&
					fieldName != "post" &&
					fieldName != "delete" &&
					fieldName != "options" &&
					fieldName != "patch" {
					continue
				}
				method := strings.ToUpper(fieldName)
				v.operations = append(v.operations,
					&operation{
						Method:    method,
						Path:      path,
						Api:       specName.ApiID,
						Version:   specName.VersionID,
						Spec:      specName.SpecID,
						Timestamp: now,
					})

			}
		}
	}
	return nil
}
