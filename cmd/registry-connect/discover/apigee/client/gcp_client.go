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
	"strings"

	"google.golang.org/api/apigee/v1"
)

func NewGCPClient() (*GCPClient, error) {
	return &GCPClient{Config.Org}, nil
}

type GCPClient struct {
	org string
}

func (c *GCPClient) Org() string {
	return strings.TrimPrefix(c.org, "organizations/")
}

func (c *GCPClient) Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("organizations/%s", c.Org())
	resp, err := apg.Organizations.Apis.List(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var proxies []*apigee.GoogleCloudApigeeV1ApiProxy
	for _, p := range resp.Proxies {
		name = fmt.Sprintf("organizations/%s/apis/%s", c.Org(), p.Name)
		proxy, err := apg.Organizations.Apis.Get(name).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

func (c *GCPClient) Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("organizations/%s", c.Org())
	resp, err := apg.Organizations.Deployments.List(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// list is minimal data, ensure detailed information
	var deps []*apigee.GoogleCloudApigeeV1Deployment
	for _, d := range resp.Deployments {
		name := fmt.Sprintf("organizations/%s/environments/%s/apis/%s/revisions/%s", c.Org(), d.Environment, d.ApiProxy, d.Revision)
		dep, err := apg.Organizations.Environments.Apis.Revisions.GetDeployments(name).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

func (c *GCPClient) EnvMap(ctx context.Context) (*EnvMap, error) {
	groups, err := c.envgroups(ctx)
	if err != nil {
		return nil, err
	}

	m := &EnvMap{
		hostnames: make(map[string][]string),
		envgroup:  make(map[string]string),
	}

	for _, group := range groups {
		name := fmt.Sprintf("organizations/%s/envgroups/%s", c.Org(), group.Name)
		attachments, err := c.attachments(ctx, name)
		if err != nil {
			return nil, err
		}

		for _, attachment := range attachments {
			for _, hostname := range group.Hostnames {
				m.hostnames[attachment.Environment] = append(m.hostnames[attachment.Environment], hostname)
				m.envgroup[hostname] = name
			}
		}
	}

	return m, nil
}

func (c *GCPClient) ProxyConsoleURL(ctx context.Context, proxy *apigee.GoogleCloudApigeeV1ApiProxy) string {
	return fmt.Sprintf("https://console.cloud.google.com/apigee/proxies/%s/overview?project=%s", proxy.Name, c.Org())
}

func (c *GCPClient) ProductConsoleURL(ctx context.Context, product *apigee.GoogleCloudApigeeV1ApiProduct) string {
	return fmt.Sprintf("https://console.cloud.google.com/apigee/apiproducts/product/%s/overview?project=%s", product.Name, c.Org())
}

func (c *GCPClient) envgroups(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1EnvironmentGroup, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("organizations/%s", c.Org())
	resp, err := apg.Organizations.Envgroups.List(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.EnvironmentGroups, nil
}

func (c *GCPClient) attachments(ctx context.Context, group string) ([]*apigee.GoogleCloudApigeeV1EnvironmentGroupAttachment, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Envgroups.Attachments.List(group).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.EnvironmentGroupAttachments, nil
}

func (c *GCPClient) Products(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("organizations/%s", c.Org())
	list, err := apg.Organizations.Apiproducts.List(name).Do()
	if err != nil {
		return nil, err
	}

	var products []*apigee.GoogleCloudApigeeV1ApiProduct
	for _, p := range list.ApiProduct {
		name := fmt.Sprintf("organizations/%s/apiproducts/%s", c.Org(), p.Name)
		product, err := apg.Organizations.Apiproducts.Get(name).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, err
}
