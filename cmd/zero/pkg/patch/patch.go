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

package patch

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v3"
)

type Header struct {
	ApiVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
}

type Metadata struct {
	Name        string            `yaml:"name" json:"name"`
	Parent      string            `yaml:"parent,omitempty" json:"parent,omitempty"`
	CreatedAt   time.Time         `yaml:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time         `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

type Patch struct {
	Header `yaml:",inline" json:",inline"`
	Data   *yaml.Node `yaml:"data" json:"data"`
}

func WrapJSON(kind, id string, bytes []byte) (*Patch, error) {
	var m yaml.Node
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return &Patch{
		Header: Header{
			ApiVersion: "zero/v1",
			Kind:       kind,
			Metadata: Metadata{
				Name: id,
			},
		},
		Data: m.Content[0],
	}, nil
}

func MarshalYAML(wrapper *Patch) ([]byte, error) {
	styleForYAML(wrapper.Data)
	bytes, err := yaml.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func MarshalJSON(wrapper *Patch) ([]byte, error) {
	bytes, err := MarshalYAML(wrapper)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	err = yaml.Unmarshal(bytes, &doc)
	if err != nil {
		return nil, err
	}
	styleForJSON(&doc)
	bytes, err = yaml.Marshal(doc.Content[0])
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func Marshal(wrapper *Patch, output string) ([]byte, error) {
	if output == "json" {
		return MarshalJSON(wrapper)
	}
	return MarshalYAML(wrapper)
}

func styleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		styleForYAML(n)
	}
}

// styleForJSON sets the style field on a tree of yaml.Nodes for JSON export.
func styleForJSON(node *yaml.Node) {
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode, yaml.MappingNode:
		node.Style = yaml.FlowStyle
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			node.Style = yaml.DoubleQuotedStyle
		default:
			node.Style = 0
		}
	case yaml.AliasNode:
	default:
	}
	for _, n := range node.Content {
		styleForJSON(n)
	}
}

func MarshalAndWrap(item interface{}, kind, name, output string) ([]byte, error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	wrapper, err := WrapJSON(kind, name, bytes)
	if err != nil {
		return nil, err
	}
	return Marshal(wrapper, output)
}
