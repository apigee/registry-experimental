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
	"strconv"
	"strings"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/cmd/registry/types"
	"github.com/apigee/registry/pkg/connection"
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
			yaml.Unmarshal(bytes, &node)

			openapi := yqQueryString(&node, "openapi")
			if openapi != nil {
				fmt.Printf("openapi %s\n", *openapi)
			}

			swagger := yqQueryString(&node, "swagger")
			if swagger != nil {
				fmt.Printf("swagger %s\n", *swagger)
			}

			title := yqQueryString(&node, "info.title")
			if title != nil {
				fmt.Printf("title %s\n", *title)
			}

			provider := yqQueryString(&node, "info.x-providerName")
			if provider != nil {
				fmt.Printf("provider %s\n", *provider)
			}
		}
		return nil
	}
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
