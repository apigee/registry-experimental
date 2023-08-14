// Copyright 2022 Google LLC.
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

	"github.com/apigee/registry/cmd/registry/patch"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/application/controller"
	"github.com/apigee/registry/pkg/application/scoring"
	"github.com/apigee/registry/pkg/application/style"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/mime"
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
	"gnostic.metrics.Complexity":                                 &metrics.Complexity{},
	"gnostic.metrics.VersionHistory":                             &metrics.VersionHistory{},
	"gnostic.metrics.Vocabulary":                                 &metrics.Vocabulary{},
	"gnostic.openapiv2.Document":                                 &oas2.Document{},
	"gnostic.openapiv3.Document":                                 &oas3.Document{},
	"google.cloud.apigeeregistry.v1.style.Lint":                  &style.Lint{},
	"google.cloud.apigeeregistry.v1.style.LintStats":             &style.LintStats{},
	"google.cloud.apigeeregistry.v1.apihub.DisplaySettings":      &apihub.DisplaySettings{},
	"google.cloud.apigeeregistry.v1.apihub.Lifecycle":            &apihub.Lifecycle{},
	"google.cloud.apigeeregistry.v1.apihub.ReferenceList":        &apihub.ReferenceList{},
	"google.cloud.apigeeregistry.v1.apihub.TaxonomyList":         &apihub.TaxonomyList{},
	"google.cloud.apigeeregistry.v1.controller.Manifest":         &controller.Manifest{},
	"google.cloud.apigeeregistry.v1.controller.Receipt":          &controller.Receipt{},
	"google.cloud.apigeeregistry.v1.scoring.Score":               &scoring.Score{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreCard":           &scoring.ScoreCard{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreCardDefinition": &scoring.ScoreCardDefinition{},
	"google.cloud.apigeeregistry.v1.scoring.ScoreDefinition":     &scoring.ScoreDefinition{},
	"google.cloud.apigeeregistry.v1.style.ConformanceReport":     &style.ConformanceReport{},
	"google.cloud.apigeeregistry.v1.style.StyleGuide":            &style.StyleGuide{},
	"google.protobuf.FileDescriptorSet":                          &descriptorpb.FileDescriptorSet{},
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
		fmt.Fprintf(w, "%s\n", contents.GetData())
		return
	}
	messageType, err := mime.MessageTypeForMimeType(contents.GetContentType())
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
	if err := patch.UnmarshalContents(contents.GetData(), contents.GetContentType(), message); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	accept := req.Header["Accept"]
	if slices.Contains(accept, "text/json") {
		// The protojson formatting doesn't include a final newline.
		fmt.Fprintln(w, protojson.Format(message))
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
	// Default return format is JSON.
	// As noted above, the protojson formatting doesn't include a final newline.
	fmt.Fprintln(w, protojson.Format(message))
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
