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

package index

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/googleapi"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Index properties of resources in the API Registry",
	}

	cmd.AddCommand(infoCommand())
	cmd.AddCommand(operationsCommand())
	cmd.AddCommand(serversCommand())
	return cmd
}

func getOrCreateDataset(ctx context.Context, client *bigquery.Client, name string) (*bigquery.Dataset, error) {
	dataset := client.Dataset(name)
	if err := dataset.Create(ctx, nil); err != nil {
		switch v := err.(type) {
		case *googleapi.Error:
			if v.Code != 409 { // already exists
				return nil, err
			}
		default:
			return nil, err
		}
	}
	return dataset, nil
}

func getOrCreateTable(ctx context.Context, dataset *bigquery.Dataset, name string, prototype interface{}) (*bigquery.Table, error) {
	table := dataset.Table(name)
	schema, err := bigquery.InferSchema(prototype)
	if err != nil {
		return nil, err
	}
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		switch v := err.(type) {
		case *googleapi.Error:
			if v.Code != 409 { // already exists
				return nil, err
			}
		default:
			return nil, err
		}
	}
	return table, nil
}
