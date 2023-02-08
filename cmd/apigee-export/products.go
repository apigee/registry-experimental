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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apigee/registry/pkg/models"
	"github.com/spf13/cobra"
	"google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v2"
)

var exportProductsCommand = &cobra.Command{
	Use:   "products ORGANIZATION [DIRECTORY]",
	Short: "Exports Apigee API Products to YAML files compatible with API Registry",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  exportProducts,
}

func exportProducts(cmd *cobra.Command, args []string) error {
	var (
		ctx = cmd.Context()
		org = args[0]
	)
	if len(args) < 2 {
		verbose = true
	}

	products, err := products(ctx, org)
	if err != nil {
		return fmt.Errorf("failed to list API products for %s: %s", org, err)
	}

	list := &models.List{
		Header: models.Header{ApiVersion: "apigeeregistry/v1"},
	}
	for _, product := range products {
		api := &models.Api{
			Header: models.Header{
				ApiVersion: "apigeeregistry/v1",
				Kind:       "API",
				Metadata: models.Metadata{
					Name: productName(org, product.Name),
					Annotations: map[string]string{
						"apigee-kind":         "product",
						"apigee-organization": org,
						"apigee-product":      fmt.Sprintf("%s/products/%s", org, product.Name),
					},
				},
			},
			Data: models.ApiData{
				DisplayName: "product: " + product.Name,
			},
		}
		list.Items = append(list.Items, api)
	}
	out, err := yaml.Marshal(list)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML for model: %s", err)
	}
	if verbose {
		fmt.Println(string(out))
	}
	// Only write the files if a directory is specified.
	if len(args) == 2 {
		filename := filepath.Join(args[1], "products.yaml")
		if err := os.WriteFile(filename, out, 0644); err != nil {
			return fmt.Errorf("failed to write YAML: %s", err)
		}
	}
	return nil
}

func products(ctx context.Context, org string) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Apiproducts.List(org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.ApiProduct, nil
}

func productName(org, name string) string {
	org = strings.TrimPrefix(org, "organizations/")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ToLower(name)
	return fmt.Sprintf("product-%s-%s", org, name)
}
