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

package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/edge"
	"google.golang.org/api/apigee/v1"
)

type EdgeClient struct {
	org string
}

func (c *EdgeClient) Org() string {
	return c.org
}

func (c *EdgeClient) Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	client, err := edge.ConfiguredClient(c.org)

	names, _, err := client.Proxies.ListNames()
	if err != nil {
		return nil, err
	}

	var proxies []*apigee.GoogleCloudApigeeV1ApiProxy
	for _, n := range names {
		proxies = append(proxies, &apigee.GoogleCloudApigeeV1ApiProxy{Name: n})
	}

	return proxies, nil
}

func (c *EdgeClient) Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	client, err := edge.ConfiguredClient(c.org)

	proxyNames, _, err := client.Proxies.ListNames()
	if err != nil {
		return nil, err
	}

	var deps []*apigee.GoogleCloudApigeeV1Deployment
	for _, n := range proxyNames {
		d, _, err := client.Proxies.GetDeployment(n)
		if err != nil {
			return nil, err
		}

		r, err := client.Proxies.GetDeployedRevision(n)
		if err != nil {
			return nil, err
		}

		for _, ed := range d.Environments {
			s := strings.Split(ed.Name, "/")
			deps = append(deps, &apigee.GoogleCloudApigeeV1Deployment{
				ApiProxy:    n,
				Environment: s[len(s)-1],
				Revision:    fmt.Sprintf("%d", r),
			})
		}
	}
	return deps, err
}

func (c *EdgeClient) EnvMap(ctx context.Context) (*EnvMap, error) {
	client, err := edge.ConfiguredClient(c.org)

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

		for _, vhName := range vhNames {
			vh, _, err := client.Environments.GetVirtualHost(envName, vhName)
			if err != nil {
				return nil, err
			}

			dedup := map[string]interface{}{}
			for _, e := range vh.HostAliases {
				dedup[e] = true
			}
			for k := range dedup {
				m.hostnames[envName] = append(m.hostnames[envName], k)
			}
		}
	}

	return m, err
}
