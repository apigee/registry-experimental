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

package proxies

import (
	"context"
	"fmt"
	"os"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/common"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	"github.com/apigee/registry/pkg/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "proxies ORGANIZATION",
		Short: "Export Apigee Proxies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			org := args[0]
			client := common.Client(org)
			return exportProxies(ctx, client)
		},
	}
	return cmd
}

func exportProxies(ctx context.Context, client common.ApigeeClient) error {
	proxies, err := client.Proxies(ctx)
	if err != nil {
		return err
	}

	var apis []interface{}
	apisByName := map[string]*encoding.Api{}
	for _, p := range proxies {
		proxy, err := client.Proxy(ctx, p.Name)
		if err != nil {
			return err
		}

		api := &encoding.Api{
			Header: encoding.Header{
				ApiVersion: encoding.RegistryV1,
				Kind:       "API",
				Metadata: encoding.Metadata{
					Name: fmt.Sprintf("%s-%s-proxy", client.Org(), proxy.Name),
					Annotations: map[string]string{
						"apigee-proxy": fmt.Sprintf("%s/apis/%s", client.Org(), proxy.Name),
					},
					Labels: map[string]string{
						"apihub-kind":          "proxy",
						"apihub-business-unit": client.Org(),
					},
				},
			},
			Data: encoding.ApiData{
				DisplayName: fmt.Sprintf("%s-%s-proxy", client.Org(), proxy.Name),
			},
		}

		for _, r := range proxy.Revision {
			v := &encoding.ApiVersion{
				Header: encoding.Header{
					ApiVersion: encoding.RegistryV1,
					Kind:       "Version",
					Metadata: encoding.Metadata{
						Name: r,
						Annotations: map[string]string{
							"apigee-proxy-revision": fmt.Sprintf("organizations/%s/apis/%s/revisions/%s", client.Org(), proxy.Name, r),
						},
					},
				},
				Data: encoding.ApiVersionData{
					DisplayName: r,
					Description: r,
				},
			}
			api.Data.ApiVersions = append(api.Data.ApiVersions, v)
		}

		rl := &apihub.ReferenceList{
			References: []*apihub.ReferenceList_Reference{{
				Id:          proxy.Name,
				DisplayName: proxy.Name + " (Apigee)",
				Uri:         client.ProxyURL(ctx, proxy),
			}},
		}
		node, err := encoding.NodeForMessage(rl)
		if err != nil {
			return err
		}

		a := &encoding.Artifact{
			Header: encoding.Header{
				ApiVersion: encoding.RegistryV1,
				Kind:       "ReferenceList",
				Metadata: encoding.Metadata{
					Name: "apihub-related",
				},
			},
			Data: *node,
		}
		api.Data.Artifacts = append(api.Data.Artifacts, a)

		apis = append(apis, api)
		apisByName[proxy.Name] = api
	}

	err = addDeployments(ctx, client, apisByName)
	if err != nil {
		return err
	}

	items := &encoding.List{
		Header: encoding.Header{ApiVersion: encoding.RegistryV1},
		Items:  apis,
	}
	return yaml.NewEncoder(os.Stdout).Encode(items)
}

func addDeployments(ctx context.Context, client common.ApigeeClient, apisByName map[string]*encoding.Api) error {
	if len(apisByName) == 0 {
		return nil
	}

	envMap, err := client.EnvMap(ctx)
	if err != nil {
		return err
	}

	deps, err := client.Deployments(ctx)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		hostnames, ok := envMap.Hostnames(dep.Environment)
		if !ok {
			log.Warnf(ctx, "Failed to find hostnames for environment %s", dep.Environment)
			continue
		}

		for _, hostname := range hostnames {
			api, ok := apisByName[dep.ApiProxy]
			if !ok {
				return fmt.Errorf("unknown proxy: %q for deployment", dep.ApiProxy)
			}

			envgroup, _ := envMap.Envgroup(hostname)
			deployment := &encoding.ApiDeployment{
				Header: encoding.Header{
					ApiVersion: encoding.RegistryV1,
					Kind:       "Deployment",
					Metadata: encoding.Metadata{
						Name: common.Label(hostname),
						Annotations: map[string]string{
							"apigee-proxy-revision": fmt.Sprintf("organizations/%s/apis/%s/revisions/%s", client.Org(), dep.ApiProxy, dep.Revision),
							"apigee-environment":    fmt.Sprintf("organizations/%s/environments/%s", client.Org(), dep.Environment),
							"apigee-envgroup":       envgroup,
						},
					},
				},
				Data: encoding.ApiDeploymentData{
					DisplayName: fmt.Sprintf("%s (%s)", dep.Environment, hostname),
					// TODO: should use proxy base path instead of name
					EndpointURI: fmt.Sprintf("https://%s/%s", hostname, dep.ApiProxy),
				},
			}

			api.Data.ApiDeployments = append(api.Data.ApiDeployments, deployment)
		}
	}
	return nil
}

/*
Example output:

apiVersion: apigeeregistry/v1
items:
  - apiVersion: apigeeregistry/v1
    kind: API
    metadata:
      name: myorg-helloworld-proxy
      labels:
        apihub-kind: proxy
        apihub-business-unit: myorg
      annotations:
        apigee-proxy: organizations/myorg/apis/helloworld
    data:
      displayName: myorg-helloworld-proxy
      deployments:
        - kind: Deployment
          metadata:
            name: prod-1-helloworld-example-com
            labels:
              apihub-gateway: apihub-google-cloud-apigee
            annotations:
              apigee-proxy-revision: organizations/myorg/apis/helloworld/revisions/1
              apigee-environment: organizations/myorg/environments/prod
              apigee-envgroup: organizations/myorg/envgroups/prod-envgroup
          data:
            displayName: prod (helloworld.example.com)
            endpointURI: helloworld.example.com
      versions:
        - kind: Version
          metadata:
            name: 1
            annotations:
              apigee-proxy-revision: organizations/myorg/apis/helloworld/revisions/1
          data:
            displayName: 1 (My First Revision)
            description: Hello World API proxy, the first revision.
      artifacts:
        - kind: ReferenceList
          metadata:
            name: apihub-related
          data:
            references:
              - id: helloworld
                displayName: helloworld (Apigee)
                uri: https://console.cloud.google.com/apigee/proxies/helloworld/overview?project=myorg

*/
