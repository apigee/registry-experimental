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

package common

import (
	"context"
	"strings"

	"github.com/apigee/registry/rpc"
	"google.golang.org/api/apigee/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type ApigeeClient interface {
	Org() string
	Proxies(context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProxy, error)
	Proxy(context.Context, string) (*apigee.GoogleCloudApigeeV1ApiProxy, error)
	Deployments(context.Context) ([]*apigee.GoogleCloudApigeeV1Deployment, error)
	EnvMap(context.Context) (*EnvMap, error)
	ProxyURL(context.Context, *apigee.GoogleCloudApigeeV1ApiProxy) string
	Products(ctx context.Context) ([]*apigee.GoogleCloudApigeeV1ApiProduct, error)
	Product(ctx context.Context, name string) (*apigee.GoogleCloudApigeeV1ApiProduct, error)
}

func Client(org string) ApigeeClient {
	// TODO: differentiate X from SaaS (and OPDK) properly
	if strings.HasPrefix(org, "organizations/") {
		return &GCPClient{org}
	} else {
		return &EdgeClient{org}
	}
}

type EnvMap struct {
	hostnames map[string][]string
	envgroup  map[string]string // only valid for X
}

func (m *EnvMap) Hostnames(env string) ([]string, bool) {
	if m.hostnames == nil {
		return nil, false
	}

	v, ok := m.hostnames[env]
	return v, ok
}

func (m *EnvMap) Envgroup(hostname string) (string, bool) {
	if m.envgroup == nil {
		return "", false
	}

	v, ok := m.envgroup[hostname]
	return v, ok
}

func Label(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(s)
}

func ArtifactNode(m *rpc.ReferenceList) (*yaml.Node, error) {
	var node *yaml.Node
	// Marshal the artifact content as JSON using the protobuf marshaller.
	s, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		Indent:          "  ",
		UseProtoNames:   false,
	}.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON with yaml.v3 so that we can re-marshal it as YAML.
	var doc yaml.Node
	err = yaml.Unmarshal([]byte(s), &doc)
	if err != nil {
		return nil, err
	}
	// The top-level node is a "document" node. We need to marshal the node below it.
	node = doc.Content[0]
	// Restyle the YAML representation so that it will be serialized with YAML defaults.
	styleForYAML(node)
	// We exclude the id and kind fields from YAML serializations.
	node = removeIdAndKind(node)
	return node, nil
}

// styleForYAML sets the style field on a tree of yaml.Nodes for YAML export.
func styleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		styleForYAML(n)
	}
}

func removeIdAndKind(node *yaml.Node) *yaml.Node {
	if node.Kind == yaml.MappingNode {
		content := make([]*yaml.Node, 0)
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			if k.Value != "id" && k.Value != "kind" {
				content = append(content, node.Content[i])
				content = append(content, node.Content[i+1])
			}
		}
		node.Content = content
	}
	return node
}
