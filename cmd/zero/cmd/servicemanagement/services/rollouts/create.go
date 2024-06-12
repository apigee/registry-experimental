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

package rollouts

import (
	"context"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/apigee/registry-experimental/cmd/zero/pkg/patch"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"

	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:  "create SERVICE CONFIG",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			serviceId := args[0]
			rolloutId := args[1]
			srv, err := servicemanagement.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			rollout := &servicemanagement.Rollout{
				RolloutId: rolloutId,
				TrafficPercentStrategy: &servicemanagement.TrafficPercentStrategy{
					Percentages: map[string]float64{
						rolloutId: 100.0,
					},
				},
			}
			item, err := srv.Services.Rollouts.
				Create(serviceId, rollout).
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
