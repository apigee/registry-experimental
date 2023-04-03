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

package backstage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apigee/registry-experimental/cmd/registry-connect/publish/backstage/encoding"
	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

const (
	apiHubTag            = "apihub"
	apiLinkFormat        = "https://pantheon.corp.google.com/apigee/hub/apis/%s/overview?project=%s"
	taxonomiesLinkFormat = "https://pantheon.corp.google.com/apigee/hub/settings/taxonomies?project=%s"
)

func Command() *cobra.Command {
	var filter string
	var cmd = &cobra.Command{
		Use:   "backstage [OUTPUT FOLDER]",
		Short: "Export APIs for a Backstage.io project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			config, err := connection.ActiveConfig()
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get config")
			}
			client, err := connection.NewRegistryClientWithSettings(ctx, config)
			if err != nil {
				return err
			}

			catalog := catalog{
				client: client,
				config: config,
				filter: filter,
				root:   args[0],
			}
			return catalog.Run(ctx)
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "filter selected apis")
	return cmd
}

func recommendedOrLatestVersion(ctx context.Context, client connection.RegistryClient, a *rpc.Api) (*rpc.ApiVersion, error) {
	n, _ := names.ParseApi(a.Name)
	versionName := n.Version("-")

	if a.RecommendedVersion != "" {
		rv, err := names.ParseVersion(a.RecommendedVersion)
		if err != nil {
			return nil, err
		}
		versionName = rv
	}

	var version *rpc.ApiVersion
	err := visitor.ListVersions(ctx, client, versionName, "", func(ctx context.Context, av *rpc.ApiVersion) error {
		version = av
		return nil
	})
	return version, err
}

func primaryOrLatestSpec(ctx context.Context, client connection.RegistryClient, av *rpc.ApiVersion) (*rpc.ApiSpec, error) {
	n, _ := names.ParseVersion(av.Name)
	specName := n.Spec("-")

	if av.PrimarySpec != "" {
		as, err := names.ParseSpec(av.PrimarySpec)
		if err != nil {
			return nil, err
		}
		specName = as
	}

	var spec *rpc.ApiSpec
	err := visitor.ListSpecs(ctx, client, specName, "", true, func(ctx context.Context, as *rpc.ApiSpec) error {
		spec = as
		return nil
	})
	return spec, err
}

func name(str string) string {
	if len(str) > 63 {
		str = str[0:63]
	}
	return str
}

func required(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

type catalog struct {
	client      connection.RegistryClient
	config      connection.Config
	filter      string
	root        string
	filesByKind map[string][]string
}

func (c *catalog) Run(ctx context.Context) error {
	c.filesByKind = map[string][]string{}

	if err := c.createGroups(ctx); err != nil {
		return err
	}
	if err := c.createAPIs(ctx); err != nil {
		return err
	}
	return c.writeCatalog()
}

func (c *catalog) createGroups(ctx context.Context) error {
	taxonomiesName, err := names.ParseArtifact(c.config.FQName("artifacts/apihub-taxonomies"))
	if err != nil {
		return err
	}
	return visitor.GetArtifact(ctx, c.client, taxonomiesName, true, func(ctx context.Context, a *rpc.Artifact) error {
		message, err := mime.MessageForMimeType(a.GetMimeType())
		if err == nil {
			err = proto.Unmarshal(a.GetContents(), message)
		}
		if err != nil {
			return err
		}
		artifactName, _ := names.ParseArtifact(a.Name)
		taxonomies := message.(*apihub.TaxonomyList)
		for _, t := range taxonomies.GetTaxonomies() {
			if t.Id == "apihub-team" {
				for _, team := range t.Elements {
					metadata := encoding.Metadata{
						Name:        name(team.Id),
						Namespace:   artifactName.ProjectID(),
						Title:       team.DisplayName,
						Description: team.Description,
						Tags:        []string{apiHubTag},
						Links: []encoding.Link{
							{
								URL:   fmt.Sprintf(taxonomiesLinkFormat, artifactName.ProjectID()),
								Title: "API Hub Taxonomies",
								// Icon: "", // requires backstage setup
							},
						},
					}
					group := encoding.Group{
						Type: "team",
					}
					if err := c.addEntity(metadata, group); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (c *catalog) createAPIs(ctx context.Context) error {
	project, err := names.ParseProject("projects/" + c.config.Project)
	if err != nil {
		return err
	}
	return visitor.ListAPIs(ctx, c.client, project.Api("-"), c.filter, func(ctx context.Context, a *rpc.Api) error {
		log.FromContext(ctx).Infof("publishing %s", a.Name)

		apiName, _ := names.ParseApi(a.Name)
		name := apiName.ApiID
		metadata := encoding.Metadata{
			Name:        name,
			Namespace:   project.ProjectID,
			Title:       a.DisplayName,
			Description: a.Description,
			Labels:      a.Labels,      // TODO: not viewable in backstage
			Annotations: a.Annotations, // TODO: not viewable in backstage
			Tags:        []string{apiHubTag},
			Links: []encoding.Link{ // TODO: not viewable in backstage
				{
					URL:   fmt.Sprintf(apiLinkFormat, apiName.ApiID, project.ProjectID),
					Title: "API Hub",
				},
			},
		}
		var style, lifecycle, definition, owner string

		owner = a.Labels["apihub-team"]
		style = strings.TrimPrefix(a.Labels["apihub-style"], "apihub-")
		lifecycle = a.Labels["apihub-lifecycle"]
		// TODO: add contact as user?
		// primaryContact = a.Labels["apihub-primary-contact"]
		// primaryContactDescription = a.Labels["apihub-primary-contact-description"]

		av, err := recommendedOrLatestVersion(ctx, c.client, a)
		if err != nil {
			return err
		}
		if av != nil {
			lifecycle = av.State

			as, err := primaryOrLatestSpec(ctx, c.client, av)
			if err != nil {
				return err
			}
			if as != nil && as.MimeType != "application/x.proto+zip" { // no binary
				// overrides API style as Backstage actually just expects spec style
				// TODO: full mime list
				if strings.Contains(as.MimeType, "openapi") || strings.Contains(as.MimeType, "yaml") {
					style = "openapi"
				} else if strings.Contains(as.MimeType, "proto") || strings.Contains(as.MimeType, "grpc") {
					style = "grpc"
				}
				definition = string(as.Contents)
			}
			// the following doesn't work to write the spec content in a separate local file,
			// perhaps related to: https://github.com/backstage/backstage/issues/14372
			// if as != nil && len(as.Contents) > 0 {
			// 	file := filepath.Join(c.root, metadata.Name+".spec")
			// 	if strings.Contains(as.MimeType, "yaml") || strings.Contains(as.MimeType, "openapi") {
			// 		definition = "$yaml: " + file
			// 	} else if strings.Contains(as.MimeType, "json") {
			// 		definition = "$json: " + file
			// 	} else if strings.Contains(as.MimeType, "text") {
			// 		definition = "$text: " + file
			// 	}
			// 	if err := os.MkdirAll(filepath.Dir(file), os.FileMode(0755)); err != nil { // rwx,rx,rx
			// 		return err
			// 	}
			// 	if err := os.WriteFile(file, as.Contents, os.FileMode(0644)); err != nil {
			// 		return err
			// 	}
			// }
		}

		api := encoding.Api{
			Type:       required(style),     // backstage well-known types: openapi, asyncapi, graphql, grpc
			Lifecycle:  required(lifecycle), // backstage well-known types: experimental, production, deprecated
			Owner:      required(owner),
			Definition: required(definition),
		}
		return c.addEntity(metadata, api)
	})
}

func (c *catalog) writeYAML(file string, data *encoding.Envelope) error {
	file = filepath.Join(c.root, file)
	if err := os.MkdirAll(filepath.Dir(file), os.FileMode(0755)); err != nil { // rwx,rx,rx
		return err
	}

	out, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644)) // rw,r,r
	if err != nil {
		return err
	}
	defer out.Close()
	enc := yaml.NewEncoder(out)
	return enc.Encode(data)
}

func (c *catalog) addEntity(metadata encoding.Metadata, entity interface{}) error {
	envelope, err := encoding.NewEnvelope(metadata, entity)
	if err != nil {
		return err
	}
	file := metadata.Name + ".yaml"
	c.filesByKind[envelope.Kind] = append(c.filesByKind[envelope.Kind], file)
	return c.writeYAML(filepath.Join(envelope.Kind+"s", file), envelope)
}

func (c *catalog) writeCatalog() error {
	subCatalogs := []string{}

	for k, fs := range c.filesByKind {
		relativeFs := []string{}
		for _, f := range fs {
			relativeFs = append(relativeFs, "./"+f)
		}
		location := encoding.Location{
			Targets: relativeFs,
		}
		pluralKinds := k + "s"
		catalog := fmt.Sprintf("All-%s.yaml", pluralKinds)
		subCatalogs = append(subCatalogs, "./"+filepath.Join(pluralKinds, catalog))
		metadata := encoding.Metadata{
			Name:        "APIHub-" + pluralKinds,
			Description: fmt.Sprintf("API Hub %s for Backstage.io", pluralKinds),
		}
		envelope, err := encoding.NewEnvelope(metadata, location)
		if err != nil {
			return err
		}
		file := filepath.Join(pluralKinds, catalog)
		if err := c.writeYAML(file, envelope); err != nil {
			return err
		}
	}

	location := encoding.Location{
		Targets: subCatalogs,
	}
	metadata := encoding.Metadata{
		Name:        "APIHub-" + c.config.Project,
		Description: fmt.Sprintf("API Hub project %s export for Backstage.io", c.config.Project),
	}
	envelope, err := encoding.NewEnvelope(metadata, location)
	if err != nil {
		return err
	}
	return c.writeYAML("APIHub-Catalog.yaml", envelope)
}
