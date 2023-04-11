// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edge

import (
	"net/url"
)

const deploymentsPath = "deployments"

// https://apidocs.apigee.com/docs/deployments/1/routes/organizations/%7Borg_name%7D/deployments/get
type OrganizationDeployments struct {
	Environments []DeploymentEnvironment `json:"environment"`
	Name         string                  `json:"name"`
}

type DeploymentEnvironment struct {
	Name       string            `json:"name,omitempty"`
	APIProxies []DeploymentProxy `json:"aPIProxy,omitempty"`
}

type DeploymentProxy struct {
	Name      string               `json:"name,omitempty"`
	Revisions []DeploymentRevision `json:"revision,omitempty"`
}

type DeploymentRevision struct {
	// Configuration interface{} `json:"configuration,omitempty"` // includeApiConfig = false, so this is not included
	Name string `json:"name,omitempty"`
	// Servers       []DeploymentServer `json:"server,omitempty"` // includeServerStatus = false, so this is not include
	State string `json:"state,omitempty"`
}

type DeploymentConfiguration struct {
	BasePath string      `json:"basePath,omitempty"`
	Steps    interface{} `json:"steps,omitempty"`
}

type DeploymentServer struct {
	Pod    Pod      `json:"pod,omitempty"`
	Status string   `json:"status,omitempty"`
	Types  []string `json:"type,omitempty"`
	UUID   string   `json:"uUID,omitempty"`
}

type Pod struct {
	Name   string `json:"name,omitempty"`
	Region string `json:"region,omitempty"`
}

// DeploymentsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apiproxies.
type DeploymentsService interface {
	OrganizationDeployments() (*OrganizationDeployments, *Response, error)
}

// DeploymentsServiceOp represents deployments
type DeploymentsServiceOp struct {
	client *EdgeClient
}

var _ DeploymentsService = &DeploymentsServiceOp{}

func (s *DeploymentsServiceOp) OrganizationDeployments() (*OrganizationDeployments, *Response, error) {
	deployments := &OrganizationDeployments{}
	q := url.Values{}
	q.Set("includeServerStatus", "false")
	q.Set("includeApiConfig", "false")
	urlString := deploymentsPath + "?" + q.Encode()
	req, e := s.client.NewRequestNoEnv("GET", urlString, deployments)
	if e != nil {
		return nil, nil, e
	}
	resp, e := s.client.Do(req, &deployments)
	if e != nil {
		return nil, resp, e
	}
	return deployments, resp, e
}
