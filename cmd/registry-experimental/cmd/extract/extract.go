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

package extract

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/cmd/registry/types"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract PATTERN",
		Short: "Extract properties from specs and artifacts stored in the registry",
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
			registryClient, err := connection.NewRegistryClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			v := &extractVisitor{
				ctx:            ctx,
				registryClient: registryClient,
			}
			return visitor.Visit(ctx, v, visitor.VisitorOptions{
				RegistryClient: registryClient,
				Pattern:        pattern,
				Filter:         filter,
			})
		},
	}
	cmd.Flags().String("filter", "", "Filter selected resources")
	cmd.Flags().Int("jobs", 10, "Number of actions to perform concurrently")
	return cmd
}

type extractVisitor struct {
	visitor.Unsupported
	ctx            context.Context
	registryClient connection.RegistryClient
}

var empty = ""

func (v *extractVisitor) SpecHandler() visitor.SpecHandler {
	return func(spec *rpc.ApiSpec) error {
		fmt.Printf("%s %s\n", spec.Name, spec.MimeType)
		err := visitor.FetchSpecContents(v.ctx, v.registryClient, spec)
		if err != nil {
			return err
		}
		bytes := spec.Contents
		if types.IsGZipCompressed(spec.MimeType) {
			bytes, err = core.GUnzippedBytes(bytes)
			if err != nil {
				return err
			}
		}
		if types.IsOpenAPIv2(spec.MimeType) || types.IsOpenAPIv3(spec.MimeType) {
			var node yaml.Node
			if err := yaml.Unmarshal(bytes, &node); err != nil {
				return err
			}

			openapi := yqQueryString(&node, "openapi")
			if openapi != nil {
				fmt.Printf("openapi %s\n", *openapi)
			}

			swagger := yqQueryString(&node, "swagger")
			if swagger != nil {
				fmt.Printf("swagger %s\n", *swagger)
			}

			description := yqQueryString(&node, "info.description")
			if description != nil {
				fmt.Printf("description %s\n", *description)
			}
			if description == nil {
				description = &empty
			}

			title := yqQueryString(&node, "info.title")
			if title != nil {
				fmt.Printf("title %s\n", *title)
			}
			if title == nil {
				title = &empty
			}

			provider := yqQueryString(&node, "info.x-providerName")
			if provider != nil {
				fmt.Printf("provider %s\n", *provider)
			}

			categories := yqQueryNode(&node, "info.x-apisguru-categories")
			if categories != nil {
				fmt.Printf("categories:\n%s\n", yqDescribe(categories))
			}

			// Set API (displayName, description) from (title, description).
			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(v.ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["openapi"] = "true"
			delete(labels, "style-openapi")
			labels["categories"] = strings.Join(yqQueryStringArray(categories), ",")
			if provider != nil {
				labels["provider"] = *provider
			}
			_, err = v.registryClient.UpdateApi(v.ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: *title,
						Description: *description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

			// Set the spec mimetype (this should not bump the revision!).
			if openapi != nil || swagger != nil {
				var compression string
				if types.IsGZipCompressed(spec.MimeType) {
					compression = "+gzip"
				}
				var mimeType string
				if openapi != nil {
					mimeType = types.OpenAPIMimeType(compression, *openapi)
				} else if swagger != nil {
					mimeType = types.OpenAPIMimeType(compression, *swagger)
				}
				specName, _ := names.ParseSpec(spec.Name)
				_, err := v.registryClient.UpdateApiSpec(v.ctx,
					&rpc.UpdateApiSpecRequest{
						ApiSpec: &rpc.ApiSpec{
							Name:     specName.String(),
							MimeType: mimeType,
						},
					},
				)
				if err != nil {
					return err
				}
			}
		}
		if types.IsDiscovery(spec.MimeType) {
			var node yaml.Node
			if err := yaml.Unmarshal(bytes, &node); err != nil {
				return err
			}
			styleForYAML(&node)
			//fmt.Printf("discovery:\n%s\n", yqDescribe(&node))

			description := yqQueryString(&node, "description")
			if description != nil {
				fmt.Printf("description %s\n", *description)
			}
			if description == nil {
				description = &empty
			}

			title := yqQueryString(&node, "canonicalName")
			if title != nil {
				fmt.Printf("title %s\n", *title)
			}
			if title == nil {
				title = &empty
			}

			provider := yqQueryString(&node, "ownerDomain")
			if provider != nil {
				fmt.Printf("provider %s\n", *provider)
			}

			// Set API (displayName, description) from (title, description).
			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(v.ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["discovery"] = "true"
			delete(labels, "style-discovery")
			if provider != nil {
				labels["provider"] = *provider
			}
			_, err = v.registryClient.UpdateApi(v.ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: *title,
						Description: *description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

		}
		if types.IsProto(spec.MimeType) {
			// create a tmp directory
			root, err := os.MkdirTemp("", "extract-protos-")
			if err != nil {
				return err
			}
			// whenever we finish, delete the tmp directory
			defer os.RemoveAll(root)
			// unzip the protos to the temp directory
			_, err = core.UnzipArchiveToPath(spec.Contents, root)
			if err != nil {
				return err
			}

			var displayName string
			var description string

			if err = filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip everything that's not a YAML file.
				if info.IsDir() || !strings.HasSuffix(filepath, ".yaml") {
					return nil
				}

				bytes, err := os.ReadFile(filepath)
				if err != nil {
					return err
				}

				sc := &ServiceConfig{}
				if err := yaml.Unmarshal(bytes, sc); err != nil {
					return err
				}

				// Skip invalid API service configurations.
				if sc.Type != "google.api.Service" || sc.Title == "" || sc.Name == "" {
					return nil
				}

				displayName = sc.Title
				description = strings.ReplaceAll(sc.Documentation.Summary, "\n", " ")

				// Skip the directory after we find an API service configuration.
				return fs.SkipDir
			}); err != nil {
				return err
			}

			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(v.ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["grpc"] = "true"
			delete(labels, "style-grpc")
			labels["provider"] = "google.com"
			_, err = v.registryClient.UpdateApi(v.ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: displayName,
						Description: description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

		}
		return nil
	}
}

// The API Service Configuration contains important API properties.
type ServiceConfig struct {
	Type          string `yaml:"type"`
	Name          string `yaml:"name"`
	Title         string `yaml:"title"`
	Documentation struct {
		Summary string `yaml:"summary"`
	} `yaml:"documentation"`
}

// styleForYAML sets the style field on a tree of yaml.Nodes for YAML export.
func styleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		styleForYAML(n)
	}
}

func yqQueryNode(node *yaml.Node, path string) *yaml.Node {
	return query(node, strings.Split(path, "."))
}

func yqQueryString(node *yaml.Node, path string) *string {
	if n := query(node, strings.Split(path, ".")); n == nil {
		return nil
	} else {
		if n.Kind == yaml.ScalarNode {
			return &n.Value
		} else {
			bytes, _ := yaml.Marshal(n)
			s := string(bytes)
			return &s
		}
	}
}

func yqQueryStringArray(node *yaml.Node) []string {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	results := []string{}
	for _, n := range node.Content {
		results = append(results, n.Value)
	}
	return results
}

func query(node *yaml.Node, path []string) *yaml.Node {
	if len(path) == 0 {
		return node
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, c := range node.Content {
			if n := query(c, path); n != nil {
				return n
			}
		}
	case yaml.SequenceNode:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return nil
		}
		return query(node.Content[index], path[1:])
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == path[0] {
				return query(node.Content[i+1], path[1:])
			}
		}
	case yaml.ScalarNode:
		return node
	case yaml.AliasNode:
		return nil
	default:
		return nil
	}
	return nil
}

func yqDescribe(node *yaml.Node) string {
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
