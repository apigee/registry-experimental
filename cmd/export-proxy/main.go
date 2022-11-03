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
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/apigee/registry/pkg/config"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"github.com/spf13/cobra"
)

type APIProxy struct {
	Name           string `xml:"name,attr"`
	DisplayName    string
	Description    string
	BasePath       string `xml:"BasePaths"`
	ProxyEndpoint  string `xml:"ProxyEndpoints>ProxyEndpoint"`
	TargetEndpoint string `xml:"TargetEndpoints>TargetEndpoint"`
}

type RouteRule struct {
	Name           string `xml:"name,attr"`
	TargetEndpoint string
}

type ProxyEndpoint struct {
	Name      string `xml:"name,attr"`
	BasePath  string `xml:"HTTPProxyConnection>BasePath"`
	RouteRule RouteRule
}

type TargetEndpoint struct {
	Name string `xml:"name,attr"`
	URL  string `xml:"HTTPTargetConnection>URL"`
}

func main() {
	cmd := &cobra.Command{
		Use:   "export-proxy",
		Short: "Exports Apigee resources to YAML files compatible with API Registry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			config, err := config.Active()
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to load active configuration")
			} else if config.Registry.Token == "" {
				log.FromContext(ctx).Fatal("Active configuration doesn't have a GCP access token")
			}

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}

			name, err := names.ParseDeployment(args[0])
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to parse deployment")
			}

			deployment, err := client.GetApiDeployment(ctx, &rpc.GetApiDeploymentRequest{
				Name: name.String(),
			})
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get deployment")
			}

			root := APIProxy{
				Name:           name.DeploymentID,
				DisplayName:    deployment.DisplayName,
				Description:    deployment.Description,
				BasePath:       "/",
				ProxyEndpoint:  "default",
				TargetEndpoint: "default",
			}

			proxy := ProxyEndpoint{
				Name:     "default",
				BasePath: "/",
				RouteRule: RouteRule{
					Name:           "default",
					TargetEndpoint: "default",
				},
			}

			target := TargetEndpoint{
				Name: "default",
				URL:  deployment.EndpointUri,
			}

			proxyZip, err := bundle(root, proxy, target)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to create bundle")
			}

			if err := createProxy(ctx, root.Name, proxyZip.Bytes(), config.Registry.Token); err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to create proxy")
			}
		},
	}

	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bundle(root APIProxy, proxy ProxyEndpoint, target TargetEndpoint) (bytes.Buffer, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	if err := write(zw, "apiproxy/"+root.Name+".xml", root); err != nil {
		return buf, err
	}

	if err := write(zw, "apiproxy/proxies/default.xml", proxy); err != nil {
		return buf, err
	}

	if err := write(zw, "apiproxy/targets/default.xml", target); err != nil {
		return buf, err
	}

	if err := zw.Close(); err != nil {
		return buf, err
	}

	return buf, nil
}

func write(zw *zip.Writer, name string, v any) error {
	out, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w, err := zw.Create(name)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	if _, err := w.Write(out); err != nil {
		return err
	}

	return nil
}

func createProxy(ctx context.Context, name string, bundle []byte, token string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://apigee.googleapis.com/v1/organizations/egansean-integrations1/apis?action=import&name=%s", name), bytes.NewReader(bundle))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Authorization", "Bearer "+token)

	_, err = http.DefaultClient.Do(req)
	return err
}
