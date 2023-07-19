// Copyright 2023 Google LLC.
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

package apigee

import (
	"context"
	"fmt"
	"os"

	apigee "github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	"github.com/apigee/registry/pkg/log"
	api "google.golang.org/api/apigee/v1"
	"gopkg.in/yaml.v3"
)

func export(ctx context.Context, client apigee.Client) error {
	log.FromContext(ctx).Infof("retrieving products")
	products, err := client.Products(ctx)
	if err != nil {
		return err
	}
	log.FromContext(ctx).Infof("%d products discovered", len(products))

	log.FromContext(ctx).Infof("retrieving proxies")
	proxies, err := client.Proxies(ctx)
	if err != nil {
		return err
	}
	log.FromContext(ctx).Infof("%d proxies discovered", len(proxies))
	proxyByName := map[string]*api.GoogleCloudApigeeV1ApiProxy{}
	for _, p := range proxies {
		proxyByName[p.Name] = p
	}

	var apis []interface{}
	for _, product := range products {
		log.FromContext(ctx).Infof("encoding product %q", product.Name)
		access := ""
		for _, a := range product.Attributes {
			if a.Name == "access" {
				access = a.Value
				break
			}
		}

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
						"apihub-target-users":  access,
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
			log.FromContext(ctx).Infof("encoding bound proxies %q", proxyNames)
			related := &apihub.ReferenceList{
				DisplayName: "Related resources",
				Description: "Links to resources in the registry.",
			}
			for _, proxyName := range proxyNames {
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

	proxyAPIs, err := addProxies(ctx, client, proxies)
	if err != nil {
		return err
	}
	apis = append(apis, proxyAPIs...)

	items := &encoding.List{
		Header: encoding.Header{ApiVersion: encoding.RegistryV1},
		Items:  apis,
	}
	log.FromContext(ctx).Infof("encoding yaml output")
	err = yaml.NewEncoder(os.Stdout).Encode(items)
	log.FromContext(ctx).Infof("products export complete")
	return err
}

func addProxies(ctx context.Context, client apigee.Client, proxies []*api.GoogleCloudApigeeV1ApiProxy) (apis []interface{}, err error) {
	apisByProxyName := map[string]*encoding.Api{}
	for _, proxy := range proxies {
		log.FromContext(ctx).Infof("encoding proxy %q", proxy.Name)
		api := &encoding.Api{
			Header: encoding.Header{
				ApiVersion: encoding.RegistryV1,
				Kind:       "API",
				Metadata: encoding.Metadata{
					Name: name(fmt.Sprintf("%s-%s-proxy", client.Org(), proxy.Name)),
					Annotations: map[string]string{
						"apigee-proxy": fmt.Sprintf("%s/apis/%s", client.Org(), proxy.Name),
					},
					Labels: map[string]string{
						"apihub-kind":          "proxy",
						"apihub-business-unit": label(client.Org()),
					},
				},
			},
			Data: encoding.ApiData{
				DisplayName: fmt.Sprintf("%s proxy: %s", client.Org(), proxy.Name),
			},
		}

		dependencies := &apihub.ReferenceList{
			DisplayName: "Apigee Dependencies",
			Description: "Links to dependant Apigee resources.",
			References: []*apihub.ReferenceList_Reference{{
				Id:          proxy.Name,
				DisplayName: proxy.Name + " (Apigee)",
				Uri:         client.ProxyConsoleURL(ctx, proxy),
			}},
		}
		node, err := encoding.NodeForMessage(dependencies)
		if err != nil {
			return nil, err
		}

		a := &encoding.Artifact{
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
		apis = append(apis, api)
		apisByProxyName[proxy.Name] = api
	}

	err = addDeployments(ctx, client, apisByProxyName)
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func addDeployments(ctx context.Context, client apigee.Client, apisByProxyName map[string]*encoding.Api) error {
	if len(apisByProxyName) == 0 {
		return nil
	}

	log.FromContext(ctx).Infof("retrieving envmap")
	envMap, err := client.EnvMap(ctx)
	if err != nil {
		return err
	}

	log.FromContext(ctx).Infof("retrieving deployments")
	deps, err := client.Deployments(ctx)
	if err != nil {
		return err
	}
	log.FromContext(ctx).Infof("%d deployments discovered", len(deps))

	for _, dep := range deps {
		hostnames, ok := envMap.Hostnames(dep.Environment)
		if !ok {
			log.FromContext(ctx).Warnf("failed to find hostnames for environment %s", dep.Environment)
			continue
		}

		for _, hostname := range hostnames {
			api, ok := apisByProxyName[dep.ApiProxy]
			if !ok {
				log.FromContext(ctx).Warnf("unknown proxy: %q for deployment: %#v", dep.ApiProxy, dep)
				continue
			}

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
	return nil
}

func boundProxies(prod *api.GoogleCloudApigeeV1ApiProduct) []string {
	proxies := prod.Proxies
	if prod.OperationGroup != nil {
		for _, oc := range prod.OperationGroup.OperationConfigs {
			if oc.ApiSource != "" {
				proxies = append(proxies, oc.ApiSource)
			}
		}
	}
	return proxies
}
