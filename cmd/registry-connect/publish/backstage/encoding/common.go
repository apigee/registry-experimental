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

package encoding

import (
	"fmt"
)

// https://backstage.io/docs/features/software-catalog/descriptor-format#overall-shape-of-an-entity

const BackstageV1alpha1 = "backstage.io/v1alpha1"

// https://backstage.io/docs/features/software-catalog/descriptor-format#common-to-all-kinds-the-envelope
type Envelope struct {
	ApiVersion string      `yaml:"apiVersion,omitempty"`
	Kind       string      `yaml:"kind,omitempty"`
	Metadata   Metadata    `yaml:"metadata,omitempty"`
	Relations  Relations   `yaml:"relations,omitempty"`
	Spec       interface{} `yaml:"spec,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/descriptor-format#common-to-all-kinds-the-metadata
type Metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Title       string            `yaml:"title,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
	Links       []Link            `yaml:"links,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/descriptor-format#links-optional
type Link struct {
	URL   string `yaml:"url,omitempty"`
	Title string `yaml:"title,omitempty"`
	Icon  string `yaml:"icon,omitempty"`
	Type  string `yaml:"type,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/descriptor-format#common-to-all-kinds-relations
type Relations struct {
	Target Reference `yaml:"target,omitempty"`
	Type   string    `yaml:"type,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/references
// [<kind>:][<namespace>/]<name>
type Reference string

func NewEnvelope(metadata Metadata, spec interface{}) (*Envelope, error) {
	var kind string
	switch t := spec.(type) {
	case Api:
		kind = "API"
	case Group:
		kind = "Group"
	case Location:
		kind = "Location"
	default:
		return nil, fmt.Errorf("invalid spec type: %#v", t)
	}
	return &Envelope{
		ApiVersion: BackstageV1alpha1,
		Kind:       kind,
		Metadata:   metadata,
		Relations:  Relations{},
		Spec:       spec,
	}, nil
}
