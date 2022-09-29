// Copyright 2022 Google LLC. All Rights Reserved.
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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/rpc"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v3"

	metrics "github.com/google/gnostic/metrics"
	oas2 "github.com/google/gnostic/openapiv2"
	oas3 "github.com/google/gnostic/openapiv3"
)

var messageTypes map[string]proto.Message = map[string]proto.Message{
	"gnostic.metrics.Complexity":                                   &metrics.Complexity{},
	"gnostic.metrics.VersionHistory":                               &metrics.VersionHistory{},
	"gnostic.metrics.Vocabulary":                                   &metrics.Vocabulary{},
	"gnostic.openapiv2.Document":                                   &oas2.Document{},
	"gnostic.openapiv3.Document":                                   &oas3.Document{},
	"google.cloud.apigeeregistry.applications.v1alpha1.Index":      &rpc.Index{},
	"google.cloud.apigeeregistry.applications.v1alpha1.Lint":       &rpc.Lint{},
	"google.cloud.apigeeregistry.applications.v1alpha1.LintStats":  &rpc.LintStats{},
	"google.cloud.apigeeregistry.applications.v1alpha1.References": &rpc.References{},
	"google.cloud.apigeeregistry.v1.apihub.DisplaySettings":        &rpc.DisplaySettings{},
	"google.cloud.apigeeregistry.v1.apihub.Lifecycle":              &rpc.Lifecycle{},
	"google.cloud.apigeeregistry.v1.apihub.ReferenceList":          &rpc.ReferenceList{},
	"google.cloud.apigeeregistry.v1.apihub.TaxonomyList":           &rpc.TaxonomyList{},
	"google.cloud.apigeeregistry.v1.controller.Manifest":           &rpc.Manifest{},
	"google.cloud.apigeeregistry.v1.controller.Receipt":            &rpc.Receipt{},
	"google.cloud.apigeeregistry.v1.scoring.Score":                 &rpc.Score{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreCard":             &rpc.ScoreCard{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreCardDefinition":   &rpc.ScoreCardDefinition{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreDefinition":       &rpc.ScoreDefinition{},
	"google.cloud.apigeeregistry.v1.style.ConformanceReport":       &rpc.ConformanceReport{},
	"google.cloud.apigeeregistry.v1.style.StyleGuide":              &rpc.StyleGuide{},
	"google.protobuf.FileDescriptorSet":                            &descriptorpb.FileDescriptorSet{},
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := connection.NewRegistryClient(ctx)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	name := strings.TrimPrefix(req.URL.Path, "/")
	contents, err := client.GetArtifactContents(ctx, &rpc.GetArtifactContentsRequest{Name: name})
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	if contents.GetContentType() == "text/plain" {
		fmt.Fprintf(w, "%s\n", string(contents.GetData()))
		return
	}
	messageType, err := core.MessageTypeForMimeType(contents.GetContentType())
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	message := messageTypes[messageType]
	if message == nil {
		err = fmt.Errorf("unsupported message type: %s", messageType)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if err := proto.Unmarshal(contents.GetData(), message); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	accept := req.Header["Accept"]
	if slices.Contains(accept, "text/json") {
		fmt.Fprint(w, protojson.Format(message))
		fmt.Fprint(w, "\n")
		return
	}
	if slices.Contains(accept, "text/yaml") {
		fmt.Fprint(w, yamlFormat(message))
		return
	}
	if slices.Contains(accept, "text/plain") {
		fmt.Fprint(w, prototext.Format(message))
		return
	}
	// default return format is JSON
	fmt.Fprint(w, protojson.Format(message))
	fmt.Fprint(w, "\n")
}

func yamlFormat(message proto.Message) string {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(protojson.Format(message)), &node)
	if err != nil {
		return err.Error()
	}
	styleForYAML(&node)
	b, err := yaml.Marshal(node.Content[0])
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func styleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		styleForYAML(n)
	}
}

func writeError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func main() {
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("%s", err)
	}
}
