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

package generate

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/cmd/registry/tasks"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
)

func openapiCommand() *cobra.Command {
	var specID string
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "Generate an OpenAPI spec from another specification format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				return fmt.Errorf("failed to get filter from flags: %s", err)
			}

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to get client: %s", err)
			}
			// Initialize task queue.
			taskQueue, wait := tasks.WorkerPoolIgnoreError(ctx, 1)
			defer wait()

			// Generate tasks.
			name := args[0]
			spec, err := names.ParseSpec(name)
			if err != nil {
				return fmt.Errorf("%q is not a valid spec name: %s", name, err)
			}

			// Iterate through a collection of specs and evaluate each.
			err = visitor.ListSpecs(ctx, client, spec, filter, false, func(ctx context.Context, spec *rpc.ApiSpec) error {
				taskQueue <- &generateOpenAPITask{
					client:    client,
					specName:  spec.Name,
					newSpecID: specID,
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to list specs: %s", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&specID, "spec-id", "generated", "ID to use for generated spec")
	return cmd
}

type generateOpenAPITask struct {
	client    connection.RegistryClient
	specName  string
	newSpecID string
}

func (task *generateOpenAPITask) String() string {
	return fmt.Sprintf("generate openapi for %s", task.specName)
}

func (task *generateOpenAPITask) Run(ctx context.Context) error {
	log.FromContext(ctx).Info(task.String())

	request := &rpc.GetApiSpecRequest{
		Name: task.specName,
	}
	spec, err := task.client.GetApiSpec(ctx, request)
	if err != nil {
		return err
	}
	data, err := visitor.GetBytesForSpec(ctx, task.client, spec)
	if err != nil {
		return err
	}
	relation := task.newSpecID
	var openapi string
	switch {
	case mime.IsOpenAPIv2(spec.GetMimeType()) || mime.IsOpenAPIv3(spec.GetMimeType()):
		return nil
	case mime.IsProto(spec.GetMimeType()) && mime.IsZipArchive(spec.GetMimeType()):
		log.FromContext(ctx).Debugf("Computing %s/specs/%s", spec.Name, relation)
		openapi, err = openAPIFromZippedProtos(ctx, spec.Name, data)
		if err != nil {
			return fmt.Errorf("error processing protos %s: %s", spec.Name, err)
		}
	case mime.IsDiscovery(spec.GetMimeType()):
		log.FromContext(ctx).Debugf("Computing %s/specs/%s", spec.Name, relation)
		openapi, err = openAPIFromDiscovery(ctx, spec.Name, data)
		if err != nil {
			return fmt.Errorf("error processing discovery %s: %s", spec.Name, err)
		}
	default:
		return fmt.Errorf("we don't know how to generate OpenAPI for %s", spec.Name)
	}

	specName, _ := names.ParseSpec(spec.GetName())
	messageData := []byte(openapi)
	messageData, err = compress.GZippedBytes(messageData)
	if err != nil {
		return fmt.Errorf("failed to compress generated OpenAPI: %s", err)
	}
	newSpec := &rpc.ApiSpec{
		Name:     specName.Version().Spec(relation).String(),
		MimeType: "application/x.openapi+gzip;version=3",
		Contents: messageData,
		Filename: "openapi.yaml",
	}
	_, err = task.client.UpdateApiSpec(ctx, &rpc.UpdateApiSpecRequest{
		ApiSpec:      newSpec,
		AllowMissing: true,
	})
	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI: %s", err)
	}
	return nil
}

// openAPIFromZippedProtos runs the OpenAPI generator and returns the results.
// This uses protoc and https://github.com/google/gnostic/tree/master/cmd/protoc-gen-openapi
// which can be installed using
//
//	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
func openAPIFromZippedProtos(ctx context.Context, name string, b []byte) (string, error) {
	// create a tmp directory
	root, err := ioutil.TempDir("", "registry-protos-")
	if err != nil {
		return "", err
	}
	// whenever we finish, delete the tmp directory
	defer os.RemoveAll(root)
	// unzip the protos to the temp directory
	_, err = compress.UnzipArchiveToPath(b, root+"/protos")
	if err != nil {
		return "", err
	}
	return generateOpenAPIForDirectory(ctx, name, root)
}

func generateOpenAPIForDirectory(ctx context.Context, name string, root string) (string, error) {
	// run protoc on all of the protos in the main directory
	protos := []string{}
	err := filepath.Walk(root+"/protos",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".proto") {
				protos = append(protos, strings.TrimPrefix(path, root+"/"))
			}
			return nil
		})
	if err != nil {
		return "", err
	}
	parts := []string{}
	parts = append(parts, protos...)
	parts = append(parts, "-I")
	parts = append(parts, "protos")
	parts = append(parts, "--openapi_out=.")
	cmd := exec.Command("protoc", parts...)
	cmd.Dir = root
	log.FromContext(ctx).Debugf("Running %+v\n", cmd)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf(
			"See %q for installation instructions",
			"https://github.com/google/gnostic/tree/main/cmd/protoc-gen-openapi")
		return "", err
	}
	log.FromContext(ctx).Debugf("protoc output: %s\n", string(data))
	// attempt to read an openapi.yaml file
	bytes, err := ioutil.ReadFile(root + "/openapi.yaml")
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// openAPIFromDiscovery runs the OpenAPI generator and returns the results.
// This uses https://github.com/LucyBot-Inc/api-spec-converter
// which can be installed using
//
//	npm install -g api-spec-converter
func openAPIFromDiscovery(ctx context.Context, name string, b []byte) (string, error) {
	// create a tmp directory
	root, err := ioutil.TempDir("", "registry-disco-")
	if err != nil {
		return "", err
	}
	// whenever we finish, delete the tmp directory
	defer os.RemoveAll(root)
	// write the spec to a file and run the converter
	err = ioutil.WriteFile(root+"/discovery.json", b, 0666)
	if err != nil {
		return "", err
	}
	args := []string{"--from", "google", "--to", "openapi_3", "discovery.json"}
	cmd := exec.Command("api-spec-converter", args...)
	cmd.Dir = root
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf(
			"See %q for installation instructions",
			"https://www.npmjs.com/package/api-spec-converter")
		return "", err
	}
	return string(data), err
}
