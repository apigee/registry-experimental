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

package yamlquery

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func QueryNode(node *yaml.Node, path string) *yaml.Node {
	return query(node, strings.Split(path, "."))
}

func QueryString(node *yaml.Node, path string) *string {
	if n := query(node, strings.Split(path, ".")); n == nil {
		return nil
	} else {
		if n.Kind == yaml.ScalarNode {
			return &n.Value
		} else {
			bytes, _ := yaml.Marshal(n)
			s := string(bytes)
			return &s
		}
	}
}

func QueryStringArray(node *yaml.Node) []string {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	results := []string{}
	for _, n := range node.Content {
		results = append(results, n.Value)
	}
	return results
}

func query(node *yaml.Node, path []string) *yaml.Node {
	if len(path) == 0 {
		return node
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, c := range node.Content {
			if n := query(c, path); n != nil {
				return n
			}
		}
	case yaml.SequenceNode:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return nil
		}
		return query(node.Content[index], path[1:])
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == path[0] {
				return query(node.Content[i+1], path[1:])
			}
		}
	case yaml.ScalarNode:
		return node
	case yaml.AliasNode:
		return nil
	default:
		return nil
	}
	return nil
}

func Describe(node *yaml.Node) string {
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
