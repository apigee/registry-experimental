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

package products

import (
	"context"
	"fmt"
	"os"
	"strings"

	apigee "github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	"github.com/apigee/registry/pkg/log"
	"github.com/spf13/cobra"
	api "google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "products ORGANIZATION",
		Short: "Export Apigee Products",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			apigee.Config.Org = args[0]
			client, err := apigee.NewClient()
			if err != nil {
				return err
			}
			return exportProducts(ctx, client)
		},
	}
	return cmd
}

func exportProducts(ctx context.Context, client apigee.Client) error {
	products, err := client.Products(ctx)
	if err != nil {
		return err
	}

	proxies, err := client.Proxies(ctx)
	if err != nil {
		return err
	}
	proxyByName := map[string]*api.GoogleCloudApigeeV1ApiProxy{}
	for _, p := range proxies {
		proxyByName[p.Name] = p
	}

	var apis []interface{}
	apisByProxy := map[string][]*encoding.Api{}
	for _, p := range products {
		product, err := client.Product(ctx, p.Name)
		if err != nil {
			return err
		}

		api := &encoding.Api{
			Header: encoding.Header{
				ApiVersion: encoding.RegistryV1,
				Kind:       "API",
				Metadata: encoding.Metadata{
					Name: fmt.Sprintf("%s-%s-product", client.Org(), product.Name),
					Annotations: map[string]string{
						"apigee-product": fmt.Sprintf("organizations/%s/apiproducts/%s", client.Org(), product.Name),
					},
					Labels: map[string]string{
						"apihub-kind":          "product",
						"apihub-business-unit": client.Org(),
						"apihub-target-users":  "internal",
					},
				},
			},
			Data: encoding.ApiData{
				DisplayName: fmt.Sprintf("%s-%s-product", client.Org(), product.Name),
				Description: fmt.Sprintf("%s API Product for internal/admin users.", product.Name),
			},
		}
		apis = append(apis, api)

		proxies := boundProxies(product)
		if len(proxies) > 0 {
			related := &apihub.ReferenceList{}
			dependencies := &apihub.ReferenceList{}
			for _, p := range proxies {
				apisByProxy[p] = append(apisByProxy[p], api)

				related.References = append(related.References, &apihub.ReferenceList_Reference{
					Id:       fmt.Sprintf("%s-%s-proxy", client.Org(), p),
					Resource: fmt.Sprintf("projects/%s/locations/global/apis/%s-%s-proxy", client.Org(), client.Org(), p),
				})

				dependencies.References = append(dependencies.References, &apihub.ReferenceList_Reference{
					Id:          p,
					DisplayName: p + " (Apigee)",
					Uri:         client.ProxyConsoleURL(ctx, proxyByName[p]),
				})
			}
			node, err := encoding.NodeForMessage(related)
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

			node, err = encoding.NodeForMessage(dependencies)
			if err != nil {
				return err
			}
			a = &encoding.Artifact{
				Header: encoding.Header{
					ApiVersion: encoding.RegistryV1,
					Kind:       "ReferenceList",
					Metadata: encoding.Metadata{
						Name: "apihub-dependencies",
					},
				},
				Data: *node,
			}
			api.Data.Artifacts = append(api.Data.Artifacts, a)
		}
	}

	err = addDeployments(ctx, client, apisByProxy)
	if err != nil {
		return err
	}

	items := &encoding.List{
		Header: encoding.Header{ApiVersion: encoding.RegistryV1},
		Items:  apis,
	}
	return yaml.NewEncoder(os.Stdout).Encode(items)
}

// product -> proxies -> deployments
func addDeployments(ctx context.Context, client apigee.Client, apisByProxy map[string][]*encoding.Api) error {
	if len(apisByProxy) == 0 {
		return nil
	}
	ps, err := client.Proxies(ctx)
	if err != nil {
		return err
	}
	proxiesByName := map[string]*api.GoogleCloudApigeeV1ApiProxy{}
	for _, p := range ps {
		proxiesByName[p.Name] = p
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
			apis, ok := apisByProxy[dep.ApiProxy]
			if !ok || len(apis) == 0 {
				log.Warnf(ctx, "unknown product: %q for deployment: %#v", dep.ApiProxy, dep)
				continue
			}

			for _, api := range apis {
				envgroup, _ := envMap.Envgroup(hostname)
				deployment := &encoding.ApiDeployment{
					Header: encoding.Header{
						ApiVersion: encoding.RegistryV1,
						Kind:       "Deployment",
						Metadata: encoding.Metadata{
							Name: label(hostname),
							Annotations: map[string]string{
								"apigee-proxy-revision": fmt.Sprintf("organizations/%s/apis/%s/revisions/%s", client.Org(), dep.ApiProxy, dep.Revision),
								"apigee-environment":    fmt.Sprintf("organizations/%s/environments/%s", client.Org(), dep.Environment),
								"apigee-envgroup":       envgroup,
							},
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName: fmt.Sprintf("%s (%s)", dep.Environment, hostname),
						// TODO - https://{org-name}-{env-name}.apigee.net/{base-path}/{resource-path} ?
						// EndpointURI: client.DeploymentEndpointURI(ctx, hostname, dep),
					},
				}
				api.Data.ApiDeployments = append(api.Data.ApiDeployments, deployment)
			}
		}
	}
	return nil
}

func boundProxies(prod *api.GoogleCloudApigeeV1ApiProduct) []string {
	proxies := prod.Proxies
	for _, oc := range prod.OperationGroup.OperationConfigs {
		if oc.ApiSource != "" {
			proxies = append(proxies, oc.ApiSource)
		}
	}
	return proxies
}

func label(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(s)
}

/*
Example output:

apiVersion: apigeeregistry/v1
items:
  - apiVersion: apigeeregistry/v1
    kind: API
    metadata:
      name: myorg-helloworld-product
      labels:
        apihub-kind: product
        apihub-target-users: internal
        apihub-business-unit: myorg
      annotations:
        apigee-product: organizations/myorg/apiproducts/helloworld
    data:
      displayName: Hello World
      description: Hello World API product for internal/admin users.
      deployments:
        - kind: Deployment
          metadata:
            name: test-helloworld-2
            labels:
              apihub-gateway: apihub-google-cloud-apigee
            annotations:
              apigee-proxy-revision: organizations/myorg/apis/helloworld/revisions/2
              apigee-environment: organizations/myorg/environments/test
          data:
            displayName: test (helloworld)
            endpointURI: helloworld-test.example.com
      artifacts:
        - kind: ReferenceList
          metadata:
            name: apihub-related
          data:
            references:
              - id: myorg-helloworld-proxy
                resource: projects/myorg/locations/global/apis/myorg-helloworld-proxy
              - id: myorg-helloworld-admin-proxy
                resource: projects/myorg/locations/global/apis/myorg-helloworld-admin-proxy
        - kind: ReferenceList
          metadata:
            name: apihub-dependencies
          data:
            references:
              - id: helloworld
                displayName: helloworld (Apigee)
                uri: https://console.cloud.google.com/apigee/proxies/helloworld/overview?project=myorg
              - id: helloworld-admin
                displayName: helloworld-admin (Apigee)
                uri: https://console.cloud.google.com/apigee/proxies/helloworld-admin/overview?project=myorg
*/
