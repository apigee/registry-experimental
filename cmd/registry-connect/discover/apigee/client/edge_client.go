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

package client

import (
	"context"
	"fmt"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/edge"
	"google.golang.org/api/apigee/v1"
)

func NewEdgeClient() (client *EdgeClient, err error) {
	client = &EdgeClient{Config.Org, nil}
	client.service, err = client.newService(context.Background())
	return client, err
}

type EdgeClient struct {
	org     string
	service *edge.EdgeClient
}

func (c *EdgeClient) newService(ctx context.Context) (*edge.EdgeClient, error) {
	opts := &edge.EdgeClientOptions{
		Debug:      Config.Debug,
		MgmtURL:    Config.MgmtURL,
		GCPManaged: false,
		Org:        c.org,
		Env:        "",
		Auth: &edge.EdgeAuth{
			SkipAuth:    false,
			Username:    Config.Username,
			Password:    Config.Password,
			MFAToken:    Config.MFAToken,
			BearerToken: Config.Token,
		},
		InsecureSkipVerify: Config.SkipVerify,
	}
	return edge.NewEdgeClient(opts)
}

func (c *EdgeClient) Org() string {
	return c.org
}

func (c *EdgeClient) Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	client := c.service

	proxyNames, _, err := client.Proxies.ListNames()
	if err != nil {
		return nil, err
	}

	var proxies []*apigee.GoogleCloudApigeeV1ApiProxy
	for _, proxyName := range proxyNames {
		proxy, _, err := client.Proxies.Get(proxyName)
		if err != nil {
			return nil, err
		}

		revisions := []string{}
		for _, r := range proxy.Revisions {
			revisions = append(revisions, r.String())
		}

		proxies = append(proxies, &apigee.GoogleCloudApigeeV1ApiProxy{
			Name:             proxy.Name,
			LatestRevisionId: proxy.Revisions[len(proxy.Revisions)-1].String(),
			Revision:         revisions,
		})
	}

	return proxies, nil
}

func (c *EdgeClient) Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	client := c.service

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
	client := c.service

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
func (c *EdgeClient) ProxyConsoleURL(ctx context.Context, proxy *apigee.GoogleCloudApigeeV1ApiProxy) string {
	return fmt.Sprintf("https://apigee.com/platform/%s/proxies/%s/overview/%s", c.Org(), proxy.Name, proxy.LatestRevisionId)
}

func (c *EdgeClient) ProductConsoleURL(ctx context.Context, product *apigee.GoogleCloudApigeeV1ApiProduct) string {
	return fmt.Sprintf("https://apigee.com/platform/%s/products/%s", c.Org(), product.Name)
}

func (c *EdgeClient) Products(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	client := c.service

	names, _, err := client.Products.ListNames()
	if err != nil {
		return nil, err
	}

	var products []*apigee.GoogleCloudApigeeV1ApiProduct
	for _, n := range names {
		p, _, err := client.Products.Get(n)
		if err != nil {
			return nil, err
		}

		var attrs []*apigee.GoogleCloudApigeeV1Attribute
		for _, a := range p.Attributes {
			attr := &apigee.GoogleCloudApigeeV1Attribute{
				Name:  a.Name,
				Value: a.Value,
			}
			attrs = append(attrs, attr)
		}

		products = append(products, &apigee.GoogleCloudApigeeV1ApiProduct{
			Name:           p.Name,
			Proxies:        p.Proxies,
			OperationGroup: &apigee.GoogleCloudApigeeV1OperationGroup{},
			Attributes:     attrs,
		})
	}

	return products, nil
}
