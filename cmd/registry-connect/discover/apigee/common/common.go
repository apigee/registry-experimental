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
	"strings"

	"google.golang.org/api/apigee/v1"
)

type ApigeeClient interface {
	Org() string
	Proxies(context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error)
	Deployments(context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error)
	EnvMap(context.Context) (*EnvMap, error)
}

func Client(org string) ApigeeClient {
	// TODO: differentiate X from SaaS
	if strings.HasPrefix(org, "organizations/") {
		return &GCPClient{org}
	} else {
		return &EdgeClient{org}
	}
}

type EnvMap struct {
	hostnames map[string][]string
	envgroup  map[string]string // only valid for X
}

func (m *EnvMap) Hostnames(env string) ([]string, bool) {
	if m.hostnames == nil {
		return nil, false
	}

	v, ok := m.hostnames[env]
	return v, ok
}

func (m *EnvMap) Envgroup(hostname string) (string, bool) {
	if m.envgroup == nil {
		return "", false
	}

	v, ok := m.envgroup[hostname]
	return v, ok
}

func Label(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(s)
}
