// Copyright 2022 Google LLC. All Rights Reserved.
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

	"github.com/apigee/registry/log"
	"github.com/apigee/registry/pkg/models"
	"github.com/spf13/cobra"
	"google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v2"
)

var exportProxiesCommand = &cobra.Command{
	Use:   "proxies ORGANIZATION [DIRECTORY]",
	Short: "Exports Apigee API Proxies to YAML files compatible with API Registry",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  exportProxies,
}

func exportProxies(cmd *cobra.Command, args []string) error {
	var (
		ctx = cmd.Context()
		org = args[0]
	)
	if len(args) < 2 {
		verbose = true
	}

	apis, err := apis(ctx, org)
	if err != nil {
		log.Fatalf(ctx, "Failed to list API proxies for %s: %s", org, err)
	}

	list := &models.List{
		Header: models.Header{ApiVersion: "apigeeregistry/v1"},
	}
	for _, proxy := range apis {
		api := &models.Api{
			Header: models.Header{
				ApiVersion: "apigeeregistry/v1",
				Kind:       "API",
				Metadata: models.Metadata{
					Name: proxyName(org, proxy.Name),
					Labels: map[string]string{
						"apigee-kind": "proxy",
					},
					Annotations: map[string]string{
						"apigee-organization": org,
						"apigee-proxy":        fmt.Sprintf("%s/apis/%s", org, proxy.Name),
					},
				},
			},
			Data: models.ApiData{
				DisplayName: "proxy: " + proxy.Name,
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
		filename := filepath.Join(args[1], "proxies.yaml")
		if err := os.WriteFile(filename, out, 0644); err != nil {
			log.Errorf(ctx, "Failed to write YAML: %s", err)
		}
	}
	return nil
}

func apis(ctx context.Context, org string) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Apis.List(org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Proxies, nil
}

func proxyName(org, name string) string {
	org = strings.TrimPrefix(org, "organizations/")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ToLower(name)
	return fmt.Sprintf("proxy-%s-%s", org, name)
}
