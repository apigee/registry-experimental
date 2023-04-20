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
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

const (
	apiLinkFormat        = "https://pantheon.corp.google.com/apigee/hub/apis/%s/overview?project=%s"
	taxonomiesLinkFormat = "https://pantheon.corp.google.com/apigee/hub/settings/taxonomies?project=%s"
)

type catalog struct {
	client         connection.RegistryClient
	config         connection.Config
	filter         string
	root           string
	entitiesByKind map[string][]*encoding.Envelope
}

func (c *catalog) Run(ctx context.Context) error {
	c.entitiesByKind = map[string][]*encoding.Envelope{}

	if err := c.createGroups(ctx); err != nil {
		return err
	}
	if err := c.createAPIs(ctx); err != nil {
		return err
	}
	return c.writeCatalog()
}

func (c *catalog) apigeeOwner() (group *encoding.Envelope, err error) {
	return c.createGroup("apg-owner", ownerName, ownerDesc)
}

func (c *catalog) createDeployment(d *rpc.ApiDeployment) (deployment *encoding.Envelope, err error) {
	var org, env, gateway, owner *encoding.Envelope
	if owner, err = c.apigeeOwner(); err != nil {
		return
	}

	var orgName, envName string
	if envLabel := d.Annotations["apigee-environment"]; envLabel != "" {
		splits := strings.Split(envLabel, "/")
		orgName = splits[1]
		envName = splits[3]
	}
	if orgName != "" && envName != "" {
		if org, err = c.addEntity(&encoding.Metadata{
			Name:        "apigee-org-" + orgName,
			Title:       "Apigee Org " + orgName,
			Description: "Apigee Org " + orgName,
		}, &encoding.Domain{
			Owner: requiredRef(owner.Reference()),
		}); err != nil {
			return
		}

		if env, err = c.addEntity(&encoding.Metadata{
			Name:        "apg-env-" + orgName + "-" + envName,
			Title:       "Apigee Env " + orgName + " " + envName,
			Description: "Apigee Env " + envName + " in Org " + orgName,
		}, &encoding.System{
			Owner:  requiredRef(owner.Reference()),
			Domain: org.Reference(),
		}); err != nil {
			return
		}

		if gateway, err = c.addEntity(&encoding.Metadata{
			Name:        "apg-gw-" + orgName,
			Title:       "Apigee Gateway " + orgName,
			Description: "Apigee Gateway in Org " + orgName + ", Env: " + envName,
		}, &encoding.Component{
			Type:      "Service",
			Lifecycle: "production",
			Owner:     requiredRef(owner.Reference()),
			System:    env.Reference(),
		}); err != nil {
			return
		}
	}

	depName, _ := names.ParseDeployment(d.Name)
	depId := depName.ApiID + "-" + depName.DeploymentID
	deployment, err = c.addEntity(&encoding.Metadata{
		Name:        "apg-dep-" + depId,
		Title:       "Apigee Deployment " + firstOf(d.DisplayName, depId),
		Description: "Apigee Deployment " + firstOf(d.DisplayName, depId) + " of API " + depName.ApiID,
		Labels:      d.Labels,
	}, &encoding.Component{
		Type:           "Service",
		Lifecycle:      required(""), // TODO: nothing to map from API Hub?
		Owner:          requiredRef(owner.Reference()),
		System:         env.Reference(),
		SubComponentOf: gateway.Reference(),
	})
	_ = gateway
	return
}

func (c *catalog) createGroups(ctx context.Context) error {
	taxonomiesName, err := names.ParseArtifact(c.config.FQName("artifacts/apihub-taxonomies"))
	if err != nil {
		return err
	}
	return visitor.GetArtifact(ctx, c.client, taxonomiesName, true, func(ctx context.Context, a *rpc.Artifact) error {
		message, err := mime.MessageForMimeType(a.GetMimeType())
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(a.GetContents(), message); err != nil {
			return err
		}
		artifactName, _ := names.ParseArtifact(a.Name)
		taxonomies := message.(*apihub.TaxonomyList)
		for _, t := range taxonomies.GetTaxonomies() {
			if t.Id == "apihub-team" {
				for _, team := range t.Elements {
					group, err := c.createGroup("apg-"+team.Id, team.DisplayName, team.Description)
					if err != nil {
						return err
					}
					group.Metadata.Links = []encoding.Link{
						{
							URL:   fmt.Sprintf(taxonomiesLinkFormat, artifactName.ProjectID()),
							Title: "API Hub Taxonomies",
						},
					}
				}
			}
		}
		return nil
	})
}

func (c *catalog) createGroup(name, title, description string) (*encoding.Envelope, error) {
	if name == "" {
		return nil, nil
	}
	return c.addEntity(&encoding.Metadata{
		Name:        name,
		Title:       firstOf(title, name),
		Description: description,
	}, &encoding.Group{
		Type: "team",
	})
}

func (c *catalog) createAPIs(ctx context.Context) error {
	project, err := names.ParseProject("projects/" + c.config.Project)
	if err != nil {
		return err
	}
	return visitor.ListAPIs(ctx, c.client, project.Api("-"), c.filter, func(ctx context.Context, a *rpc.Api) error {
		log.FromContext(ctx).Infof("publishing %s", a.Name)

		var specContents string
		style := strings.TrimPrefix(a.Labels["apihub-style"], "apihub-")
		if style == "" {
			if _, ok := a.Annotations["apigee-proxy"]; ok {
				style = "apigee-proxy"
			} else if _, ok := a.Annotations["apigee-product"]; ok {
				style = "apigee-product"
			}
		}
		lifecycle := a.Labels["apihub-lifecycle"]

		primaryContact, err := c.createGroup(a.Labels["apihub-primary-contact"], a.Labels["apihub-primary-contact"], a.Labels["apihub-primary-contact-description"])
		if err != nil {
			return err
		}

		// TODO: denormalize or take one? one for now.
		av, err := recommendedOrLatestVersion(ctx, c.client, a)
		if err != nil {
			return err
		}
		if av != nil {
			if av.State != "" {
				lifecycle = av.State
			}

			var specs []*rpc.ApiSpec
			vName, _ := names.ParseVersion(av.Name)
			err = visitor.ListSpecs(ctx, c.client, vName.Spec("-"), "", true, func(ctx context.Context, as *rpc.ApiSpec) error {
				specs = append(specs, as)
				return nil
			})
			if err != nil {
				return err
			}

			// TODO: denormalize or take one? one for now.
			var as *rpc.ApiSpec
			for _, s := range specs { // take primary
				if av.PrimarySpec == s.Name {
					as = s
				}
			}
			if as == nil { // or take first
				if len(specs) > 0 {
					as = specs[0]
				}
			}

			if as != nil && as.MimeType != "application/x.proto+zip" { // no binary
				// Backstage well-known types: openapi, asyncapi, graphql, grpc
				if strings.Contains(as.MimeType, "openapi") || strings.Contains(as.MimeType, "yaml") {
					style = "openapi"
				} else if strings.Contains(as.MimeType, "proto") || strings.Contains(as.MimeType, "grpc") {
					style = "grpc"
				} else if strings.Contains(as.MimeType, "asyncapi") {
					style = "asyncapi"
				} else if strings.Contains(as.MimeType, "graphql") {
					style = "graphql"
				}
				specContents = string(as.Contents)
			}

			apiName, _ := names.ParseApi(a.Name)
			apiHubLinks := []encoding.Link{{
				URL:   fmt.Sprintf(apiLinkFormat, apiName.ApiID, c.config.Project),
				Title: "API Hub",
			}}
			api, err := c.addEntity(&encoding.Metadata{
				Name:        "apg-" + apiName.ApiID,
				Title:       "Apigee " + firstOf(a.DisplayName, apiName.ApiID),
				Description: firstOf(a.Description, a.DisplayName),
				Labels:      a.Labels, // note: labels and links are not viewable in default backstage API plugin
				Links:       apiHubLinks,
			}, &encoding.Api{
				Type:       required(style),
				Lifecycle:  required(lifecycle), // backstage well-known types: experimental, production, deprecated
				Owner:      requiredRef(primaryContact.Reference()),
				Definition: required(specContents),
			})
			if err != nil {
				return err
			}

			err = visitor.ListDeployments(ctx, c.client, apiName.Deployment("-"), "", func(ctx context.Context, d *rpc.ApiDeployment) error {
				env, err := c.createDeployment(d)
				dep := env.Spec.(*encoding.Component)
				dep.ProvidesApis = append(dep.ProvidesApis, api.Reference())
				env.Metadata.Links = append(env.Metadata.Links, apiHubLinks...)
				return err
			})
			return err
		}
		return err
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

func (c *catalog) defaultNamespace() string {
	return c.config.Project
}

// will create entity if Reference (kind, namespace, name) doesn't exist,
// otherwise will return pointer to existing entity
// namespace will default, an "apiHub" tag will be added
// name and namespace may be modified to be valid
func (c *catalog) addEntity(metadata *encoding.Metadata, spec encoding.Spec) (*encoding.Envelope, error) {
	if metadata.Namespace == "" {
		metadata.Namespace = c.defaultNamespace()
	}
	metadata.Tags = append(metadata.Tags, "apihub")
	envelope, err := encoding.NewEnvelope(metadata, spec)
	if err != nil {
		return nil, err
	}
	if env := c.findEntity(envelope.Kind, envelope.Metadata.Namespace, envelope.Metadata.Name); env != nil {
		return env, nil
	}
	c.entitiesByKind[envelope.Kind] = append(c.entitiesByKind[envelope.Kind], envelope)
	return envelope, nil
}

func (c *catalog) findEntity(kind, namespace, name string) *encoding.Envelope {
	for _, env := range c.entitiesByKind[kind] {
		if env.Metadata.Namespace == encoding.SafeName(namespace) && env.Metadata.Name == encoding.SafeName(name) {
			return env
		}
	}
	return nil
}

func (c *catalog) writeCatalog() error {
	subCatalogs := []string{}

	for k, entities := range c.entitiesByKind {
		files := []string{}
		pluralKind := strings.ToLower(k) + "s"
		for _, entity := range entities {
			safeName := encoding.SafeName(string(entity.Reference()))
			fileName := strings.ToLower(safeName) + ".yaml"
			if err := c.writeYAML(filepath.Join(pluralKind, fileName), entity); err != nil {
				return err
			}
			files = append(files, "./"+fileName)
		}

		kindCatalog, err := encoding.NewEnvelope(&encoding.Metadata{
			Name:        "apihub-" + pluralKind,
			Description: fmt.Sprintf("API Hub %s for Backstage.io", pluralKind),
		}, &encoding.Location{
			Targets: files,
		})
		if err != nil {
			return err
		}

		fileName := filepath.Join(pluralKind, fmt.Sprintf("all-%s.yaml", pluralKind))
		if err := c.writeYAML(fileName, kindCatalog); err != nil {
			return err
		}
		subCatalogs = append(subCatalogs, "./"+fileName)
	}

	catalog, err := encoding.NewEnvelope(&encoding.Metadata{
		Name:        "apihub-" + strings.ToLower(c.config.Project),
		Description: fmt.Sprintf("API Hub project %s export for Backstage.io", c.config.Project),
	}, &encoding.Location{
		Targets: subCatalogs,
	})
	if err != nil {
		return err
	}
	return c.writeYAML("apihub-catalog.yaml", catalog)
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

func required(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func requiredRef(ref encoding.Reference) encoding.Reference {
	if ref == "" {
		return "unknown"
	}
	return ref
}

func firstOf(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}
