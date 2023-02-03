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
	"github.com/apigee/registry/pkg/names"
	registry_rpc "github.com/apigee/registry/rpc"

	discovery_v1 "github.com/google/gnostic/discovery"
)

func isDiscovery(s *registry_rpc.ApiSpec) bool {
	return strings.Contains(s.MimeType, "discovery")
}

func newDocumentsForDiscovery(spec *registry_rpc.ApiSpec, document *discovery_v1.Document) ([]*models.Document, error) {
	s, err := names.ParseSpec(spec.Name)
	if err != nil {
		return nil, err
	}
	text := fmt.Sprintf("%s\n%s\n%s\n%s",
		document.OwnerName,
		document.Title,
		document.Version,
		document.Description)
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
