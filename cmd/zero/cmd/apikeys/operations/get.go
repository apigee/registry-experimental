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

package operations

import (
	"context"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/option"

	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:  "get OPERATION",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			operation := args[0]
			srv, err := apikeys.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			item, err := srv.Operations.
				Get(operation).
				Do()
			if err != nil {
				return err
			}
			bytes, err := patch.MarshalAndWrap(item, "Operation", item.Name, output)
			if err != nil {
				return err
			}
			if _, err := cmd.OutOrStdout().Write(bytes); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of: (yaml, json).")
	return cmd
}
