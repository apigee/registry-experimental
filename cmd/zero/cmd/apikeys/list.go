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

package apikeys

import (
	"context"
	"fmt"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/option"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var output string
	var producerProject string
	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			srv, err := apikeys.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			var cursor string
			c := 0
			for i := 0; true; i++ {
				items, err := apikeys.NewProjectsLocationsKeysService(srv).
					List(fmt.Sprintf("projects/%s/locations/global", producerProject)).
					PageToken(cursor).
					Do()
				if err != nil {
					return err
				}
				for _, item := range items.Keys {
					bytes, err := patch.MarshalAndWrap(item, "ApiKey", item.Name, output)
					if err != nil {
						return err
					}
					if c > 0 && output != "json" {
						fmt.Printf("---\n")
					}
					c++
					if _, err := cmd.OutOrStdout().Write(bytes); err != nil {
						return err
					}
				}
				cursor = items.NextPageToken
				if cursor == "" {
					break
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of: (yaml, json).")
	cmd.Flags().StringVarP(&producerProject, "project", "p", "", "Producer project.")
	return cmd
}
