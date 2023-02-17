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

const deploymentsPath = "deployments"

// https://apidocs.apigee.com/docs/deployments/1/routes/organizations/%7Borg_name%7D/deployments/get
type OrganizationDeployments struct {
	Environments []Environment `json:"environment"`
	Name         string        `json:"name"`
	State        string        `json:"state"`
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
	req, e := s.client.NewRequestNoEnv("GET", deploymentsPath, deployments)
	if e != nil {
		return nil, nil, e
	}
	var names []string
	resp, e := s.client.Do(req, &names)
	if e != nil {
		return nil, resp, e
	}
	return deployments, resp, e
}
