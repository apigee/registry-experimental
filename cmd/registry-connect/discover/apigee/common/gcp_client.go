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

	"google.golang.org/api/apigee/v1"
)

type GCPClient struct {
	org string
}

func (c *GCPClient) Org() string {
	return c.org
}

func (c *GCPClient) Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Apis.List(c.org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Proxies, nil
}

func (c *GCPClient) Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Deployments.List(c.org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Deployments, nil
}

func (c *GCPClient) EnvMap(ctx context.Context) (*EnvMap, error) {
	return nil, nil
}
