// Copyright 2023 Google LLC.
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

package yamlquery_test

import (
	"testing"

	"github.com/apigee/registry-experimental/pkg/yamlquery"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestQueryString(t *testing.T) {

	testYaml := `
sub:
  list:
  - one
  - two
  map:
    foo: bar
    bar: baz
key: val
anchored: &anchor
  foo: bar
aliased: *anchor`

	node := new(yaml.Node)
	if err := yaml.Unmarshal([]byte(testYaml), node); err != nil {
		t.Fatal(err)
	}

	sp := func(v string) *string { return &v }

	tests := []struct {
		path string
		want *string
	}{
		{"sub.list.0", sp("one")},
		{"sub.list.1", sp("two")},
		{"sub.list.3", nil},
		{"sub.list.not_int", nil},
		{"sub.list", sp("- one\n- two\n")},
		{"sub.map.foo", sp("bar")},
		{"sub.map.bar", sp("baz")},
		{"sub.map.missing", nil},
		{"sub.map", sp("foo: bar\nbar: baz\n")},
		{"key", sp("val")},
		{"key.is.this.right", nil},
		{"anchored", sp("&anchor\nfoo: bar\n")},
		{"aliased", sp("*anchor\n")},
		{"aliased.foo", nil},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {

			got := yamlquery.QueryString(node, test.path)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestQueryStringArray(t *testing.T) {

	testYaml := `
sub:
  list:
  - one
  - two`

	node := new(yaml.Node)
	if err := yaml.Unmarshal([]byte(testYaml), node); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
		want []string
	}{
		{"sub", nil},
		{"sub.list", []string{"one", "two"}},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			node := yamlquery.QueryNode(node, test.path)
			got := yamlquery.QueryStringArray(node)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDescribe(t *testing.T) {

	testYaml := `sub:
    list:
        - one
        - two
    map:
        foo: bar
        bar: baz
key: val
anchored: &anchor
    foo: bar
aliased: *anchor
`

	node := new(yaml.Node)
	if err := yaml.Unmarshal([]byte(testYaml), node); err != nil {
		t.Fatal(err)
	}

	got := yamlquery.Describe(node)
	if diff := cmp.Diff(testYaml, got); diff != "" {
		t.Errorf("Unexpected diff (-want +got):\n%s", diff)
	}
}
