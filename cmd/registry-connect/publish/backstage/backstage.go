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

package backstage

import (
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var filter string
	var cmd = &cobra.Command{
		Use:   "backstage [OUTPUT FOLDER]",
		Short: "Export APIs for a Backstage.io project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			config, err := connection.ActiveConfig()
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get config")
			}
			client, err := connection.NewRegistryClientWithSettings(ctx, config)
			if err != nil {
				return err
			}

			catalog := catalog{
				client: client,
				config: config,
				filter: filter,
				root:   args[0],
			}
			return catalog.Run(ctx)
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "filter selected apis")
	return cmd
}
