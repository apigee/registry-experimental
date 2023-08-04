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
	"github.com/apigee/registry-experimental/cmd/registry-bigquery/common"
	"github.com/apigee/registry-experimental/pkg/yamlquery"
	"github.com/apigee/registry/cmd/registry/patch"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func linksCommand() *cobra.Command {
	var filter string
	var project string
	var dataset string
	var batchSize int
	cmd := &cobra.Command{
		Use:   "links PATTERN",
		Short: "Build a BigQuery index of links between resources",
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
			ds, err := common.GetOrCreateDataset(ctx, client, dataset)
			if err != nil {
				return err
			}
			table, err := common.GetOrCreateTable(ctx, ds, "operations", operation{})
			if err != nil {
				return err
			}
			v := &linksVisitor{
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
			log.Printf("uploading %d links", len(v.links))
			for start := 0; start < len(v.links); start += batchSize {
				log.Printf("%d", start)
				end := min(start+batchSize, len(v.links))
				if err := u.Put(ctx, v.links[start:end]); err != nil {
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

type link struct {
	Path      string
	Method    string
	Api       string
	Version   string
	Spec      string
	Timestamp time.Time
}

type linksVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
	links          []*link
}

func (v *linksVisitor) ArtifactHandler() visitor.ArtifactHandler {
	return func(ctx context.Context, message *rpc.Artifact) error {
		fmt.Printf("%s\n", message.Name)
		artifactName, err := names.ParseArtifact(message.Name)
		if err != nil {
			return err
		}
		kind := mime.KindForMimeType(message.MimeType)
		if kind != "ReferenceList" {
			return nil // skip it
		}
		m := &apihub.ReferenceList{}
		err = visitor.FetchArtifactContents(ctx, v.registryClient, message)
		if err != nil {
			return err
		}
		if err := patch.UnmarshalContents(message.GetContents(), message.GetMimeType(), m); err != nil {
			return err
		}
		fmt.Printf("  %s\n", artifactName.ApiID())
		for _, link := range m.References {
			if link.Resource != "" {
				n, err := names.ParseApi(link.Resource)
				if err != nil {
					continue
				}
				fmt.Printf("  -->%s\n", n.ApiID)
			}
		}
		return nil
	}
}

func (v *linksVisitor) getOpenAPIOperations(specName names.Spec, b []byte) error {
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
				v.links = append(v.links,
					&link{
						Method:    method,
						Path:      path,
						Api:       specName.ApiID,
						Version:   specName.VersionID,
						Spec:      specName.SpecID,
						Timestamp: common.Now,
					})
			}
		}
	}
	return nil
}
