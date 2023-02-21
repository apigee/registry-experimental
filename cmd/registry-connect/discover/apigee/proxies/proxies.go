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

package proxies

import (
	"context"
	"fmt"
	"os"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/common"
	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/edge"
	"github.com/apigee/registry/cmd/registry/patch"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/models"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "proxies ORGANIZATION",
		Short: "Export Apigee Proxies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			org := args[0]
			client := common.Client(org)
			return exportProxies(ctx, client)
		},
	}
	cmd.Flags().BoolVar(&edge.Debug, "debug", false, "debug")
	return cmd
}

func exportProxies(ctx context.Context, client common.ApigeeClient) error {
	proxies, err := client.Proxies(ctx)
	if err != nil {
		return err
	}

	apis := map[string]*models.Api{}
	for _, proxy := range proxies {
		api := &models.Api{
			Header: models.Header{
				ApiVersion: patch.RegistryV1,
				Kind:       "API",
				Metadata: models.Metadata{
					Name: common.Label(proxy.Name),
					Annotations: map[string]string{
						"apigee-proxy": fmt.Sprintf("organizations/%s/apis/%s", client.Org(), proxy.Name),
					},
					Labels: map[string]string{
						"apihub-kind":          "proxy",
						"apihub-business-unit": client.Org(),
					},
				},
			},
			Data: models.ApiData{
				DisplayName: proxy.Name,
			},
		}

		for _, r := range proxy.Revision {
			v := &models.ApiVersion{
				Header: models.Header{
					ApiVersion: patch.RegistryV1,
					Kind:       "Version",
					Metadata: models.Metadata{
						Name: r,
						Annotations: map[string]string{
							"apigee-proxy-revision": r,
						},
					},
				},
				Data: models.ApiVersionData{
					DisplayName: r,
					Description: r,
				},
			}
			api.Data.ApiVersions = append(api.Data.ApiVersions, v)
		}

		// TODO: SaaS, OPDK?
		// Apigee X
		proxyURL := fmt.Sprintf("https://console.cloud.google.com/apigee/proxies/%s/overview?project=%s", proxy.Name, client.Org())
		rl := &rpc.ReferenceList{
			References: []*rpc.ReferenceList_Reference{{
				Id:          proxy.Name,
				DisplayName: proxy.Name + " (Apigee)",
				Uri:         proxyURL,
			}},
		}
		node, err := artifactNode(rl)
		if err != nil {
			return err
		}

		a := &models.Artifact{
			Header: models.Header{
				ApiVersion: patch.RegistryV1,
				Kind:       "ReferenceList",
				Metadata: models.Metadata{
					Name: "apihub-related",
				},
			},
			Data: *node,
		}
		api.Data.Artifacts = append(api.Data.Artifacts, a)

		apis[proxy.Name] = api
	}

	err = addDeployments(ctx, client, apis)
	if err != nil {
		return err
	}

	items := &struct {
		ApiVersion string
		Items      interface{}
	}{
		ApiVersion: patch.RegistryV1,
		Items:      apis,
	}
	return yaml.NewEncoder(os.Stdout).Encode(items)
}

func addDeployments(ctx context.Context, client common.ApigeeClient, apis map[string]*models.Api) error {
	env, err := client.EnvMap(ctx)
	if err != nil {
		return err
	}

	deps, err := client.Deployments(ctx)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		hostnames, ok := env.Hostnames(dep.Environment)
		if !ok {
			log.Warnf(ctx, "Failed to find hostnames for environment %s", dep.Environment)
			continue
		}

		for _, hostname := range hostnames {
			api, ok := apis[dep.ApiProxy]
			if !ok {
				return fmt.Errorf("unknown proxy: %q for deployment", dep.ApiProxy)
			}

			envgroup, _ := env.Envgroup(hostname)
			deployment := &models.ApiDeployment{
				Header: models.Header{
					ApiVersion: patch.RegistryV1,
					Kind:       "Deployment",
					Metadata: models.Metadata{
						Name: common.Label(hostname),
						Annotations: map[string]string{
							"apigee-proxy-revision": fmt.Sprintf("%s/apis/%s/revisions/%s", client.Org(), dep.ApiProxy, dep.Revision),
							"apigee-environment":    fmt.Sprintf("%s/environments/%s", client.Org(), dep.Environment),
							"apigee-envgroup":       envgroup,
						},
					},
				},
				Data: models.ApiDeploymentData{
					DisplayName: fmt.Sprintf("%s (%s)", dep.Environment, hostname),
					// TODO: should use proxy base path instead of name
					EndpointURI: fmt.Sprintf("https://%s/%s", hostname, dep.ApiProxy),
				},
			}

			api.Data.ApiDeployments = append(api.Data.ApiDeployments, deployment)
		}
	}
	return nil
}

func artifactNode(m *rpc.ReferenceList) (*yaml.Node, error) {
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

/*
Example output:

apiVersion: apigeeregistry/v1
items:
  - apiVersion: apigeeregistry/v1
    kind: API
    metadata:
      name: myorg-helloworld-proxy
      labels:
        apihub-kind: proxy
        apihub-business-unit: myorg
      annotations:
        apigee-proxy: organizations/myorg/apis/helloworld
    data:
      displayName: myorg-helloworld-proxy
      deployments:
        - kind: Deployment
          metadata:
            name: prod-1-helloworld-example-com
            labels:
              apihub-gateway: apihub-google-cloud-apigee
            annotations:
              apigee-proxy-revision: organizations/myorg/apis/helloworld/revisions/1
              apigee-environment: organizations/myorg/environments/prod
              apigee-envgroup: organizations/myorg/envgroups/prod-envgroup
          data:
            displayName: prod (helloworld.example.com)
            endpointURI: helloworld.example.com
      versions:
        - kind: Version
          metadata:
            name: 1
            annotations:
              apigee-proxy-revision: organizations/myorg/apis/helloworld/revisions/1
          data:
            displayName: 1 (My First Revision)
            description: Hello World API proxy, the first revision.
      artifacts:
        - kind: ReferenceList
          metadata:
            name: apihub-related
          data:
            references:
              - id: helloworld
                displayName: helloworld (Apigee)
                uri: https://console.cloud.google.com/apigee/proxies/helloworld/overview?project=myorg

*/
