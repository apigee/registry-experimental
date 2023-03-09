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
	"regexp"
	"strings"

	apigee "github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	"github.com/apigee/registry/pkg/log"
	"github.com/spf13/cobra"
	api "google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v3"
)

var project string // TODO: remove when a relative ReferenceList_Reference.Resource works in Hub

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
	cmd.Flags().StringVarP(&project, "project", "", "", "hub project id (temporary)")
	_ = cmd.MarkFlagRequired("project")
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
	apisByProxyName := map[string][]*encoding.Api{}
	for _, product := range products {
		api := &encoding.Api{
			Header: encoding.Header{
				ApiVersion: encoding.RegistryV1,
				Kind:       "API",
				Metadata: encoding.Metadata{
					Name: name(fmt.Sprintf("%s-%s-product", client.Org(), product.Name)),
					Annotations: map[string]string{
						"apigee-product": fmt.Sprintf("organizations/%s/apiproducts/%s", client.Org(), product.Name),
					},
					Labels: map[string]string{
						"apihub-kind":          "product",
						"apihub-business-unit": label(client.Org()),
						"apihub-target-users":  "internal",
					},
				},
			},
			Data: encoding.ApiData{
				DisplayName: fmt.Sprintf("%s product: %s", client.Org(), product.Name),
				Description: fmt.Sprintf("%s API Product for internal/admin users.", product.Name),
			},
		}
		apis = append(apis, api)

		dependencies := &apihub.ReferenceList{
			DisplayName: "Apigee Dependencies",
			Description: "Links to dependant Apigee resources.",
		}
		dependencies.References = append(dependencies.References, &apihub.ReferenceList_Reference{
			Id:          product.Name,
			DisplayName: product.Name + " (Apigee)",
			Uri:         client.ProductConsoleURL(ctx, product),
		})

		proxyNames := boundProxies(product)
		if len(proxyNames) > 0 {
			related := &apihub.ReferenceList{
				DisplayName: "Related resources",
				Description: "Links to resources in the registry.",
			}
			for _, proxyName := range proxyNames {
				apisByProxyName[proxyName] = append(apisByProxyName[proxyName], api)

				related.References = append(related.References, &apihub.ReferenceList_Reference{
					Id:          fmt.Sprintf("%s-%s-proxy", client.Org(), proxyName),
					DisplayName: fmt.Sprintf("%s proxy: %s", client.Org(), proxyName),
					Resource:    fmt.Sprintf("projects/%s/locations/global/apis/%s-%s-proxy", project, client.Org(), proxyName),
				})

				proxy := proxyByName[proxyName]
				if proxy == nil {
					log.FromContext(ctx).Warnf("proxy %q bound but not found", proxyName)
					continue
				}
				dependencies.References = append(dependencies.References, &apihub.ReferenceList_Reference{
					Id:          proxyName,
					DisplayName: proxyName + " (Apigee)",
					Uri:         client.ProxyConsoleURL(ctx, proxy),
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

	err = addDeployments(ctx, client, apisByProxyName)
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
func addDeployments(ctx context.Context, client apigee.Client, apisByProxyName map[string][]*encoding.Api) error {
	if len(apisByProxyName) == 0 {
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
			apis, ok := apisByProxyName[dep.ApiProxy]
			if !ok || len(apis) == 0 {
				log.Warnf(ctx, "Unknown proxy: %q for deployment: %#v", dep.ApiProxy, dep)
				continue
			}

			for _, api := range apis {
				envgroup, _ := envMap.Envgroup(hostname)
				deployment := &encoding.ApiDeployment{
					Header: encoding.Header{
						ApiVersion: encoding.RegistryV1,
						Kind:       "Deployment",
						Metadata: encoding.Metadata{
							Name: name(hostname),
							Annotations: map[string]string{
								"apigee-proxy-revision": fmt.Sprintf("organizations/%s/apis/%s/revisions/%s", client.Org(), dep.ApiProxy, dep.Revision),
								"apigee-environment":    fmt.Sprintf("organizations/%s/environments/%s", client.Org(), dep.Environment),
								"apigee-envgroup":       envgroup,
							},
							Labels: map[string]string{
								"apihub-gateway": "apihub-google-cloud-apigee",
							},
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName: fmt.Sprintf("%s (%s)", dep.Environment, hostname),
						EndpointURI: hostname, // TODO: full resource path?
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
	return strings.ToLower(regexp.MustCompile(`([^A-Za-z0-9-_]+)`).ReplaceAllString(s, "-"))
}

func name(s string) string {
	return strings.ToLower(regexp.MustCompile(`([^A-Za-z0-9-]+)`).ReplaceAllString(s, "-"))
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
