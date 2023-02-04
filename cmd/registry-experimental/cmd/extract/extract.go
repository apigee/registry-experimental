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
	"log"

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
		log.Printf("%s %s", spec.Name, spec.MimeType)
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
			fmt.Printf("%+v\n", node)
		}
		return nil
	}
}
