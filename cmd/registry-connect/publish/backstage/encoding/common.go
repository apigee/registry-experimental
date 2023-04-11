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
	"regexp"
	"strings"
)

// https://backstage.io/docs/features/software-catalog/descriptor-format#overall-shape-of-an-entity

const BackstageV1alpha1 = "backstage.io/v1alpha1"

type Spec interface{}

// https://backstage.io/docs/features/software-catalog/descriptor-format#common-to-all-kinds-the-envelope
type Envelope struct {
	ApiVersion string     `yaml:"apiVersion,omitempty"`
	Kind       string     `yaml:"kind,omitempty"`
	Metadata   *Metadata  `yaml:"metadata,omitempty"`
	Relations  []Relation `yaml:"relations,omitempty"`
	Spec       Spec       `yaml:"spec,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/references
// format: [<kind>:][<namespace>/]<name>
func (e *Envelope) Reference() Reference {
	if e == nil {
		return ""
	}
	var kind, namespace string
	if e.Kind != "" {
		kind = e.Kind + ":"
	}
	if e.Metadata.Namespace != "" {
		namespace = e.Metadata.Namespace + "/"
	}
	return Reference(kind + namespace + e.Metadata.Name)
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
type Relation struct {
	Target Reference `yaml:"target,omitempty"`
	Type   string    `yaml:"type,omitempty"`
}

// https://backstage.io/docs/features/software-catalog/references
// [<kind>:][<namespace>/]<name>
type Reference string

// will automatically fix name and namespace using SafeName
func NewEnvelope(metadata *Metadata, spec Spec) (*Envelope, error) {
	metadata.Name = SafeName(metadata.Name)
	metadata.Namespace = SafeName(metadata.Namespace)
	kind, err := Kind(spec)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		ApiVersion: BackstageV1alpha1,
		Kind:       kind,
		Metadata:   metadata,
		Relations:  []Relation{},
		Spec:       spec,
	}, nil
}

func Kind(spec Spec) (string, error) {
	switch t := spec.(type) {
	case *Api:
		return "API", nil
	case *Component:
		return "Component", nil
	case *Domain:
		return "Domain", nil
	case *Group:
		return "Group", nil
	case *Location:
		return "Location", nil
	case *System:
		return "System", nil
	case *User:
		return "User", nil
	default:
		return "", fmt.Errorf("invalid spec type: %#v", t)
	}
}

// Strings of length at least 1, and at most 63
// Must consist of sequences of [a-z0-9A-Z] possibly separated by one of [-_.]
func SafeName(str string) string {
	str = regexp.MustCompile(`[^a-z0-9A-Z-.]+`).ReplaceAllString(str, "_")
	if len(str) > 63 {
		str = str[0:63]
	}
	str = strings.TrimRight(str, "-._")
	return str
}
