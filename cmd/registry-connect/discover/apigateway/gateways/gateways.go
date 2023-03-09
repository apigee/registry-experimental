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

package gateways

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apigee/registry/pkg/encoding"
	"github.com/spf13/cobra"
)

const weNeedToReadAPIs = false

func Command() *cobra.Command {
	var output string
	var cmd = &cobra.Command{
		Use:   "gateways",
		Short: "Export API Gateway Gateways",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			project, err := getProject()
			if err != nil {
				return err
			}
			if weNeedToReadAPIs {
				apis, err := fetchApis()
				if err != nil {
					return err
				}
				for i, api := range apis {
					fmt.Printf("%d %+v\n", i, api)
				}
			}
			gateways, err := fetchGateways()
			if err != nil {
				return err
			}
			for _, gateway := range gateways {
				fmt.Printf("exporting %s\n", gateway.Name)
				config, err := fetchConfig(gateway.ApiConfig)
				if err != nil {
					return err
				}
				err = writeRegistryYAML(output, project, &gateway, config)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&output, "output", "gateway-export", "output directory")

	return cmd
}

func getProject() (string, error) {
	buf := &bytes.Buffer{}
	command := "gcloud config get project"
	cmdargs := strings.Split(command, " ")
	c := exec.Command(cmdargs[0], cmdargs[1:]...)
	c.Stdout = buf
	if err := c.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func fetchApis() ([]GatewayApi, error) {
	buf := &bytes.Buffer{}
	command := "gcloud api-gateway apis list --format json"
	cmdargs := strings.Split(command, " ")
	c := exec.Command(cmdargs[0], cmdargs[1:]...)
	c.Stdout = buf
	if err := c.Run(); err != nil {
		return nil, err
	}
	var apis []GatewayApi
	if err := json.Unmarshal(buf.Bytes(), &apis); err != nil {
		return nil, err
	}
	return apis, nil
}

func fetchGateways() ([]Gateway, error) {
	buf := &bytes.Buffer{}
	command := "gcloud api-gateway gateways list --format json"
	cmdargs := strings.Split(command, " ")
	c := exec.Command(cmdargs[0], cmdargs[1:]...)
	c.Stdout = buf
	if err := c.Run(); err != nil {
		return nil, err
	}
	var gateways []Gateway
	if err := json.Unmarshal(buf.Bytes(), &gateways); err != nil {
		return nil, err
	}
	return gateways, nil
}

func fetchConfig(name string) (*GatewayConfig, error) {
	buf := &bytes.Buffer{}
	command := "gcloud api-gateway api-configs describe --view FULL --format json " + name
	cmdargs := strings.Split(command, " ")
	c := exec.Command(cmdargs[0], cmdargs[1:]...)
	c.Stdout = buf
	if err := c.Run(); err != nil {
		return nil, err
	}
	var config GatewayConfig
	if err := json.Unmarshal(buf.Bytes(), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type GatewayApi struct {
	CreateTime     string `json:"createTime"`
	DisplayName    string `json:"displayName"`
	ManagedService string `json:"managedService"`
	Name           string `json:"name"`
	State          string `json:"state"`
	UpdateTime     string `json:"updateTime"`
}

type Gateway struct {
	ApiConfig       string `json:"apiConfig"`
	CreateTime      string `json:"createTime"`
	DefaultHostname string `json:"defaultHostname"`
	DisplayName     string `json:"displayName"`
	Name            string `json:"name"`
	State           string `json:"state"`
	UpdateTime      string `json:"updateTime"`
}

type GatewayConfig struct {
	CreateTime            string                     `json:"createTime"`
	DisplayName           string                     `json:"displayName"`
	GatewayServiceAccount string                     `json:"gatewayServiceAccount"`
	GrpcServices          []map[string][]GatewayFile `json:"grpcServices"`
	ManagedServiceConfigs []*GatewayFile             `json:"managedServiceConfigs"`
	OpenAPIDocuments      []map[string]GatewayFile   `json:"openapiDocuments"`
	Name                  string                     `json:"name"`
	ServiceConfigId       string                     `json:"serviceConfigId"`
	State                 string                     `json:"state"`
	UpdateTime            string                     `json:"updateTime"`
}

type GatewayFile struct {
	Contents string `json:"contents"`
	Path     string `json:"path"`
}

func writeRegistryYAML(output, project string, gateway *Gateway, config *GatewayConfig) error {
	apiId := filepath.Base(gateway.Name)
	err := os.MkdirAll(filepath.Join(output, project+"-"+apiId), 0777)
	if err != nil {
		return err
	}
	description := "Exported from API Gateway"

	specname := "openapi"
	filename := "openapi.yaml"
	mimetype := "application/x.openapi+gzip;version=2.0"
	if config.GrpcServices != nil {
		specname = "protos"
		filename = "protos.zip"
		mimetype = "application/x.proto+zip"
	}
	a := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: project + "-" + apiId,
				Labels: map[string]string{
					"provider": strings.ReplaceAll(project, ".", "-"),
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName: project + " gateway " + apiId,
			Description: description,
			ApiVersions: []*encoding.ApiVersion{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "v1",
						},
					},
					Data: encoding.ApiVersionData{
						DisplayName: "v1",
						ApiSpecs: []*encoding.ApiSpec{
							{
								Header: encoding.Header{
									Metadata: encoding.Metadata{
										Name: specname,
									},
								},
								Data: encoding.ApiSpecData{
									FileName: filename,
									MimeType: mimetype,
								},
							},
						},
					},
				},
			},
			ApiDeployments: []*encoding.ApiDeployment{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "gateway",
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName: "gateway",
						EndpointURI: "https://" + gateway.DefaultHostname,
					},
				},
			},
		},
	}
	b, err := encoding.EncodeYAML(a)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(filepath.Join(output, project+"-"+apiId, "info.yaml"), b, 0666); err != nil {
		return err
	}
	for _, openapi := range config.OpenAPIDocuments {
		f := openapi["document"]
		contents, err := base64.StdEncoding.DecodeString(f.Contents)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(output, project+"-"+apiId, f.Path), contents, 0666); err != nil {
			return err
		}
	}
	for _, grpcService := range config.GrpcServices {
		files := grpcService["source"]
		for _, f := range files {
			contents, err := base64.StdEncoding.DecodeString(f.Contents)
			if err != nil {
				return err
			}
			if err = os.WriteFile(filepath.Join(output, project+"-"+apiId, f.Path), contents, 0666); err != nil {
				return err
			}
		}
	}
	for _, f := range config.ManagedServiceConfigs {
		contents, err := base64.StdEncoding.DecodeString(f.Contents)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(output, project+"-"+apiId, f.Path), contents, 0666); err != nil {
			return err
		}
	}
	return nil
}
