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

package configs

import (
	"context"
	"fmt"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:  "list SERVICE",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return nil
			}
			service := args[0]
			srv, err := servicemanagement.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return nil
			}
			var cursor string
			c := 0
			for i := 0; true; i++ {
				items, err := srv.Services.Configs.
					List(service).
					PageToken(cursor).
					Do()
				if err != nil {
					return nil
				}
				for _, item := range items.ServiceConfigs {
					bytes, err := patch.MarshalAndWrap(item, "ServiceConfig", item.Name, output)
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
	return cmd
}
