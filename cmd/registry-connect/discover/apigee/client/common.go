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
	"strings"

	"google.golang.org/api/apigee/v1"
)

var Config = struct {
	Org        string
	Username   string
	Password   string
	MFAToken   string
	Debug      bool
	SkipVerify bool
	MgmtURL    string
	OPDK       bool
	Edge       bool
}{}

func NewClient() (Client, error) {
	if Config.OPDK || Config.Edge {
		return NewEdgeClient()
	}

	if !strings.HasPrefix(Config.Org, "organizations/") {
		Config.Org = "organizations/" + Config.Org
	}
	return NewGCPClient()
}

type Client interface {
	Org() string
	Proxies(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error)
	ProxyConsoleURL(ctx context.Context, proxy *apigee.GoogleCloudApigeeV1ApiProxy) string
	ProductConsoleURL(ctx context.Context, product *apigee.GoogleCloudApigeeV1ApiProduct) string
	Deployments(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error)
	EnvMap(ctx context.Context) (*EnvMap, error)
	Products(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error)
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
