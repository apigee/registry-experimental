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
	"github.com/spf13/cobra"
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/option"
)

func createCmd() *cobra.Command {
	var output string
	var producerProject string
	cmd := &cobra.Command{
		Use:  "create KEYID SERVICE",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			apikey := &apikeys.V2Key{
				Restrictions: &apikeys.V2Restrictions{
					ApiTargets: []*apikeys.V2ApiTarget{
						{Service: args[1]},
					},
				},
			}
			srv, err := apikeys.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			parent := fmt.Sprintf("projects/%s/locations/global", producerProject)
			keyId := args[0]
			item, err := apikeys.NewProjectsLocationsKeysService(srv).
				Create(parent, apikey).
				KeyId(keyId).
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
	cmd.Flags().StringVarP(&producerProject, "project", "p", "", "Producer project.")
	return cmd
}
