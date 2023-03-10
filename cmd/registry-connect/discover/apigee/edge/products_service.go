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

import "path"

const productsPath = "apiproducts"

// https://apidocs.apigee.com/docs/api-products/1/routes/organizations/%7Borg_name%7D/apiproducts/get

type ApiProduct struct {
	APIResources   []string    `json:"apiResources"`
	ApprovalType   string      `json:"approvalType"`
	CreatedAt      Timestamp   `json:"createdAt,omitempty"`
	CreatedBy      string      `json:"createdBy,omitempty"`
	Description    string      `json:"description"`
	DisplayName    string      `json:"displayName"`
	Environments   []string    `json:"environments"`
	LastModifiedAt Timestamp   `json:"lastModifiedAt,omitempty"`
	LastModifiedBy string      `json:"lastModifiedBy,omitempty"`
	Name           string      `json:"name"`
	Proxies        []string    `json:"proxies"`
	Scopes         []string    `json:"scopes"`
	Attributes     []Attribute `json:"attributes"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with APIProducts.
type ProductsService interface {
	ListNames() ([]string, *Response, error)
	Get(name string) (*ApiProduct, *Response, error)
}

// ProductsServiceOp represents deployments
type ProductsServiceOp struct {
	client *EdgeClient
}

var _ ProductsService = &ProductsServiceOp{}

func (s *ProductsServiceOp) ListNames() ([]string, *Response, error) {
	req, e := s.client.NewRequestNoEnv("GET", productsPath, nil)
	if e != nil {
		return nil, nil, e
	}
	var names []string
	resp, e := s.client.Do(req, &names)
	if e != nil {
		return nil, resp, e
	}
	return names, resp, e
}

func (s *ProductsServiceOp) Get(name string) (*ApiProduct, *Response, error) {
	urlPath := path.Join(productsPath, name)
	req, e := s.client.NewRequestNoEnv("GET", urlPath, nil)
	if e != nil {
		return nil, nil, e
	}
	product := &ApiProduct{}
	resp, e := s.client.Do(req, product)
	return product, resp, e
}
