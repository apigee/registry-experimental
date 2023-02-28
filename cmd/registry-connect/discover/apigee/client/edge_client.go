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

package client

import (
	"context"
	"fmt"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/edge"
	"google.golang.org/api/apigee/v1"
)

func NewEdgeClient() (*EdgeClient, error) {
	return &EdgeClient{Config.Org}, nil
}

type EdgeClient struct {
	org string
}

func (c *EdgeClient) newService(ctx context.Context) (*edge.EdgeClient, error) {
	opts := &edge.EdgeClientOptions{
		Debug:      Config.Debug,
		MgmtURL:    Config.MgmtURL,
		GCPManaged: false,
		Org:        c.org,
		Env:        "",
		Auth: &edge.EdgeAuth{
			SkipAuth: false,
			Username: Config.Username,
			Password: Config.Password,
			MFAToken: Config.MFAToken,
		},
		InsecureSkipVerify: Config.SkipVerify,
	}
	return edge.NewEdgeClient(opts)
}

func (c *EdgeClient) Org() string {
	return c.org
}

func (c *EdgeClient) Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	names, _, err := client.Proxies.ListNames()
	if err != nil {
		return nil, err
	}

	var proxies []*apigee.GoogleCloudApigeeV1ApiProxy
	for _, n := range names {
		p, _, err := client.Proxies.Get(n)
		if err != nil {
			return nil, err
		}

		proxies = append(proxies, &apigee.GoogleCloudApigeeV1ApiProxy{
			Name:             n,
			LatestRevisionId: p.Revisions[len(p.Revisions)-1].String(),
		})
	}

	return proxies, nil
}

func (c *EdgeClient) Proxy(ctx context.Context, name string) (*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	p, _, err := client.Proxies.Get(name)
	if err != nil {
		return nil, err
	}

	revisions := []string{}
	for _, r := range p.Revisions {
		revisions = append(revisions, r.String())
	}

	proxy := &apigee.GoogleCloudApigeeV1ApiProxy{
		Name:             p.Name,
		LatestRevisionId: revisions[len(revisions)-1],
		Revision:         revisions,
	}

	return proxy, nil
}

func (c *EdgeClient) Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	var deps []*apigee.GoogleCloudApigeeV1Deployment
	ods, _, err := client.Deployments.OrganizationDeployments()
	if err != nil {
		return nil, err
	}

	for _, e := range ods.Environments {
		for _, p := range e.APIProxies {
			for _, r := range p.Revisions {
				dep := &apigee.GoogleCloudApigeeV1Deployment{
					ApiProxy:    p.Name,
					Environment: e.Name,
					Revision:    r.Name,
					State:       r.State,
				}
				deps = append(deps, dep)
			}
		}
	}
	return deps, err
}

func (c *EdgeClient) EnvMap(ctx context.Context) (*EnvMap, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	envNames, _, err := client.Environments.ListNames()
	if err != nil {
		return nil, err
	}

	m := &EnvMap{
		hostnames: make(map[string][]string),
	}

	for _, envName := range envNames {
		vhNames, _, err := client.Environments.ListVirtualHosts(envName)
		if err != nil {
			return nil, err
		}

		uniqueHostnames := map[string]bool{}
		for _, vhName := range vhNames {
			vh, _, err := client.Environments.GetVirtualHost(envName, vhName)
			if err != nil {
				return nil, err
			}
			for _, e := range vh.HostAliases {
				uniqueHostnames[e] = true
			}
		}
		for k := range uniqueHostnames {
			m.hostnames[envName] = append(m.hostnames[envName], k)
		}
	}

	return m, err
}

// TODO: Won't work with OPDK
func (c *EdgeClient) ProxyURL(ctx context.Context, proxy *apigee.GoogleCloudApigeeV1ApiProxy) string {
	return fmt.Sprintf("https://apigee.com/platform/%s/proxies/%s/overview/%s", c.org, proxy.Name, proxy.LatestRevisionId)
}

func (c *EdgeClient) Products(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	names, _, err := client.Products.ListNames()
	if err != nil {
		return nil, err
	}

	var products []*apigee.GoogleCloudApigeeV1ApiProduct
	for _, n := range names {
		products = append(products, &apigee.GoogleCloudApigeeV1ApiProduct{
			Name: n,
		})
	}

	return products, nil
}

func (c *EdgeClient) Product(ctx context.Context, name string) (*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	client, err := c.newService(ctx)
	if err != nil {
		return nil, err
	}

	p, _, err := client.Products.Get(name)
	if err != nil {
		return nil, err
	}

	product := &apigee.GoogleCloudApigeeV1ApiProduct{
		Name:           p.Name,
		Proxies:        p.Proxies,
		OperationGroup: &apigee.GoogleCloudApigeeV1OperationGroup{},
	}

	return product, nil
}
