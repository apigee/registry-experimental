// Copyright 2021 Google LLC. All Rights Reserved.
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
	"fmt"
	"strings"

	"github.com/apigee/registry-experimental/server/search/internal/storage/models"
	registry_rpc "github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	openapi_v3 "github.com/google/gnostic/openapiv3"
)

func isOpenAPIv2(s *registry_rpc.ApiSpec) bool {
	return strings.Contains(s.MimeType, "openapi") && strings.Contains(s.Name, "swagger.yaml")
}

func newDocumentsForOpenAPIv2(spec *registry_rpc.ApiSpec, document *openapi_v2.Document) ([]*models.Document, error) {
	s, err := names.ParseSpec(spec.Name)
	if err != nil {
		return nil, err
	}
	text := fmt.Sprintf("%s\n%s\n%s",
		document.Info.Title,
		document.Info.Version,
		document.Info.Description)
	docs := []*models.Document{
		(&models.Document{
			Key:       spec.Name,
			Name:      spec.Name,
			Fragment:  "",
			Kind:      specEntityName,
			Field:     "",
			ProjectID: s.ProjectID,
			Vector:    models.TSVector{RawText: text, Weight: models.WeightA},
		}).Escape(),
	}
	return docs, nil
}

func isOpenAPIv3(s *registry_rpc.ApiSpec) bool {
	return strings.Contains(s.MimeType, "openapi") && strings.Contains(s.Name, "openapi.yaml")
}

func newDocumentsForOpenAPIv3(spec *registry_rpc.ApiSpec, document *openapi_v3.Document) ([]*models.Document, error) {
	s, err := names.ParseSpec(spec.Name)
	if err != nil {
		return nil, err
	}
	text := fmt.Sprintf("%s\n%s\n%s",
		document.Info.Title,
		document.Info.Version,
		document.Info.Description)
	docs := []*models.Document{
		(&models.Document{
			Key:       spec.Name,
			Name:      spec.Name,
			Fragment:  "",
			Kind:      specEntityName,
			Field:     "",
			ProjectID: s.ProjectID,
			Vector:    models.TSVector{RawText: text, Weight: models.WeightA},
		}).Escape(),
	}
	return docs, nil
}
