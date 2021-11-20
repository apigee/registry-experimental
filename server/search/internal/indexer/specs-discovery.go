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
	//	docs = appendAllDiscoveryParameters(docs, spec, document.GetParameters())
	//	docs = appendAllDiscoverySchemas(docs, spec, document.GetSchemas())
	//	docs = appendAllDiscoveryResources(docs, "", spec, document.GetResources())
	return docs, nil
}

func appendAllDiscoveryParameters(docs []*models.Document, spec *registry_rpc.ApiSpec, params *discovery_v1.Parameters) []*models.Document {
	base := "parameters/"
	for _, m := range params.GetAdditionalProperties() {
		path := base + m.GetName()
		description := m.GetValue().GetDescription()
		docs = append(docs, newDocumentForSpecPath(spec, models.FieldParameters, models.WeightD, path, description))
	}
	return docs
}

func appendAllDiscoverySchemas(docs []*models.Document, spec *registry_rpc.ApiSpec, schemas *discovery_v1.Schemas) []*models.Document {
	base := "schemas/"
	for _, m := range schemas.GetAdditionalProperties() {
		path := base + m.GetName()
		s := m.GetValue()
		description := s.GetDescription()
		docs = append(docs, newDocumentForSpecPath(spec, models.FieldSchemas, models.WeightC, path, description))
		docs = appendAllDiscoverySchemaProperties(docs, path+"/", spec, s.GetProperties())
	}
	return docs
}

func appendAllDiscoverySchemaProperties(docs []*models.Document, base string, spec *registry_rpc.ApiSpec, properties *discovery_v1.Schemas) []*models.Document {
	base += "properties/"
	for _, m := range properties.GetAdditionalProperties() {
		path := base + m.GetName()
		p := m.GetValue()
		description := p.GetDescription()
		docs = append(docs, newDocumentForSpecPath(spec, models.FieldParameters, models.WeightC, path, description))

		path += "/enums/"
		for i, e := range p.GetEnum() {
			if i >= len(p.GetEnumDescriptions()) {
				break
			}

			description = p.GetEnumDescriptions()[i]
			docs = append(docs, newDocumentForSpecPath(spec, models.FieldParameters, models.WeightC, path+e, description))
		}
	}
	return docs
}

func appendAllDiscoveryResources(docs []*models.Document, base string, spec *registry_rpc.ApiSpec, resources *discovery_v1.Resources) []*models.Document {
	base += "resources/"
	for _, r := range resources.GetAdditionalProperties() {
		path := base + r.GetName() + "/"
		docs = appendAllDiscoveryMethods(docs, path, spec, r.GetValue().GetMethods())
		docs = appendAllDiscoveryResources(docs, path, spec, r.GetValue().GetResources())
	}
	return docs
}

func appendAllDiscoveryMethods(docs []*models.Document, base string, spec *registry_rpc.ApiSpec, methods *discovery_v1.Methods) []*models.Document {
	base += "methods/"
	for _, m := range methods.GetAdditionalProperties() {
		path := base + m.GetName()
		description := m.GetValue().GetDescription()
		docs = append(docs, newDocumentForSpecPath(spec, models.FieldMethods, models.WeightB, path, description))
	}
	return docs
}
