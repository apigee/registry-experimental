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

package servicecontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"google.golang.org/api/option"
	"google.golang.org/api/servicecontrol/v1"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func checkCmd() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:  "check",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			srv, err := servicecontrol.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}
			now := time.Now()
			timestamp := now.Format(time.RFC3339)

			request := &servicecontrol.CheckRequest{
				Operation: &servicecontrol.Operation{
					OperationId:   uuid.New().String(),
					OperationName: "/hello",
					//ConsumerId:    "project:" + producerProject,
					ConsumerId: "api_key:" + apiKey,
					StartTime:  timestamp,
					Labels: map[string]string{
						"cloud.googleapis.com/service":             serviceName,
						"serviceruntime.googleapis.com/api_method": "1.hello_nbuv3ljuva_uw_a_run_app.Hello",
						"servicecontrol.googleapis.com/caller_ip":  "172.125.77.209",
						"servicecontrol.googleapis.com/user_agent": "ESPv2",
					},
				},
			}
			result, err := srv.Services.Check(serviceName, request).Do()
			if err != nil {
				return err
			}
			bytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return &json.UnsupportedValueError{}
			}
			fmt.Printf("%s", string(bytes))
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of: (yaml, json).")
	return cmd
}
