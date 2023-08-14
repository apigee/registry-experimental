// Copyright 2021 Google LLC.
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

package indexer

import (
	"github.com/apigee/registry-experimental/server/search/internal/storage/models"
	registry_rpc "github.com/apigee/registry/rpc"
	discovery_v1 "github.com/google/gnostic/discovery"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	openapi_v3 "github.com/google/gnostic/openapiv3"
)

const specEntityName = "Spec"

func NewDocumentsForSpec(spec *registry_rpc.ApiSpec, contents []byte) ([]*models.Document, error) {
	switch {
	case isDiscovery(spec):
		document, err := discovery_v1.ParseDocument(contents)
		if err != nil {
			return nil, err
		}
		return newDocumentsForDiscovery(spec, document)
	case isOpenAPIv2(spec):
		document, err := openapi_v2.ParseDocument(contents)
		if err != nil {
			return nil, err
		}
		return newDocumentsForOpenAPIv2(spec, document)
	case isOpenAPIv3(spec):
		document, err := openapi_v3.ParseDocument(contents)
		if err != nil {
			return nil, err
		}
		return newDocumentsForOpenAPIv3(spec, document)
	}
	return nil, nil
}
