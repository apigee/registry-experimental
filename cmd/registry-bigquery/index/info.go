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

func infoCommand() *cobra.Command {
	var filter string
	var project string
	var dataset string
	var step int
	cmd := &cobra.Command{
		Use:   "info PATTERN",
		Short: "Build a BigQuery index of API information",
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
			v := &infoVisitor{
				registryClient: registryClient,
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
			table, err := getOrCreateTable(ctx, ds, "info", info{})
			if err != nil {
				return err
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
			log.Printf("uploading %d info sets", len(v.infos))
			for start := 0; start < len(v.infos); start += step {
				log.Printf("%d", start)
				end := min(start+step, len(v.infos))
				if err := u.Put(ctx, v.infos[start:end]); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "Filter selected resources")
	cmd.Flags().StringVar(&project, "project", "", "Project to use for BigQuery upload (defaults to registry project)")
	cmd.Flags().StringVar(&dataset, "dataset", "registry", "BigQuery dataset name")
	cmd.Flags().IntVar(&step, "step", 10000, "Step size to use when uploading operations to BigQuery")
	return cmd
}

type info struct {
	Title       string
	Description string
	Api         string
	Version     string
	Spec        string
	Timestamp   time.Time
}

type infoVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
	infos          []*info
}

func (v *infoVisitor) SpecHandler() visitor.SpecHandler {
	return func(ctx context.Context, message *rpc.ApiSpec) error {
		fmt.Printf("%s\n", message.Name)
		specName, err := names.ParseSpec(message.Name)
		if err != nil {
			return err
		}
		return visitor.GetSpec(ctx, v.registryClient, specName, true,
			func(ctx context.Context, spec *rpc.ApiSpec) error {
				if mime.IsOpenAPIv2(spec.MimeType) || mime.IsOpenAPIv3(spec.MimeType) {
					err := v.getOpenAPIInfo(specName, spec.Contents)
					if err != nil {
						return err
					}
				}
				return nil
			})
	}
}

func (v *infoVisitor) getOpenAPIInfo(specName names.Spec, b []byte) error {
	var doc yaml.Node
	err := yaml.Unmarshal(b, &doc)
	if err != nil {
		return err
	}
	titlep := yamlquery.QueryString(&doc, "info.title")
	var title string
	if titlep != nil {
		title = *titlep
	}
	descriptionp := yamlquery.QueryString(&doc, "info.description")
	var description string
	if descriptionp != nil {
		description = *descriptionp
	}
	if title != "" || description != "" {
		v.infos = append(v.infos,
			&info{
				Title:       title,
				Description: description,
				Api:         specName.ApiID,
				Version:     specName.VersionID,
				Spec:        specName.SpecID,
				Timestamp:   now,
			})
	}
	return nil
}
