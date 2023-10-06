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

package match

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	"google.golang.org/api/iterator"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	var filter string
	var project string
	var dataset string
	var batchSize int
	cmd := &cobra.Command{
		Use:   "match PATTERN",
		Short: "Match API specs with a BigQuery index of API information",
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
			v := &matchVisitor{
				registryClient: registryClient,
				bigQueryClient: client,
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

			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "Filter selected resources")
	cmd.Flags().StringVar(&project, "project", "", "Project to use for BigQuery upload (defaults to registry project)")
	cmd.Flags().StringVar(&dataset, "dataset", "registry", "BigQuery dataset name")
	cmd.Flags().IntVar(&batchSize, "batch-size", 10000, "Batch size to use when uploading records to BigQuery")
	return cmd
}

type matchVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
	bigQueryClient *bigquery.Client
}

func (v *matchVisitor) SpecHandler() visitor.SpecHandler {
	return func(ctx context.Context, message *rpc.ApiSpec) error {
		fmt.Printf("MATCHING %s\n", message.Name)
		specName, err := names.ParseSpec(message.Name)
		if err != nil {
			return err
		}
		return visitor.GetSpec(ctx, v.registryClient, specName, true,
			func(ctx context.Context, spec *rpc.ApiSpec) error {
				if mime.IsOpenAPIv2(spec.MimeType) || mime.IsOpenAPIv3(spec.MimeType) {
					err := v.matchOpenAPI(ctx, specName, spec.Contents)
					if err != nil {
						return err
					}
				}
				return nil
			})
	}
}

func (v *matchVisitor) matchOpenAPI(ctx context.Context, specName names.Spec, b []byte) error {
	var doc yaml.Node
	err := yaml.Unmarshal(b, &doc)
	if err != nil {
		return err
	}
	operations := make([]*operation, 0)
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
				operations = append(operations,
					&operation{
						Method:  method,
						Path:    path,
						Api:     specName.ApiID,
						Version: specName.VersionID,
						Spec:    specName.SpecID,
					})
			}
		}
	}

	counts := make(map[string]int)
	total := len(operations)
	fmt.Printf("%d total operations\n", total)
	for i, op := range operations {
		fmt.Printf("%d: %s %s\n", i, op.Method, op.Path)
		pattern := strings.ReplaceAll(op.Path, "*", "%")
		query := fmt.Sprintf(
			`SELECT * FROM registry.operations WHERE path like "%s" and method = "%s"`,
			pattern,
			op.Method)
		q := v.bigQueryClient.Query(query)
		it, err := q.Read(ctx)
		if err != nil {
			return err
		}
		for {
			var match operation
			err = it.Next(&match)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			name := fmt.Sprintf("apis/%s/versions/%s/specs/%s", match.Api, match.Version, match.Spec)
			counts[name]++
		}
	}

	apis := make([]string, 0)
	for k := range counts {
		apis = append(apis, k)
	}
	sort.Slice(apis, func(i int, j int) bool {
		api_i := apis[i]
		api_j := apis[j]
		count_i := counts[api_i]
		count_j := counts[api_j]
		if count_i > count_j {
			return true
		} else if count_i < count_j {
			return false
		}
		// apis with equal counts are alphabetized
		return api_i < api_j
	})
	fmt.Println("")
	// print the array backwards so the best match is last
	for i := len(apis) - 1; i >= 0; i-- {
		api := apis[i]
		fmt.Printf("%0.5f\t%s\n", 1.0*float32(counts[api])/float32(total), api)
	}
	// if we didn't match anything, we are finished
	if len(apis) == 0 {
		return nil
	}
	// assume the last match is the one we want to save
	// create a link to from the traffic signal to the reference
	// the traffic is the api of the starting specName
	trafficApi := specName.Api()
	// the reference is the last match
	lastApi := apis[len(apis)-1]
	enrolledApi := specName.Project().Api(lastApi)

	log.Printf("LINKING %s to %s", trafficApi, enrolledApi)

	return nil
}

type operation struct {
	Path      string
	Method    string
	Api       string
	Version   string
	Spec      string
	Timestamp time.Time
}
