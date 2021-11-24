// Copyright 2021 Google LLC. All Rights Reserved.
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

	"github.com/apex/log"
	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"github.com/spf13/cobra"
)

func openapiCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Generate an OpenAPI spec for a protocol buffer API specification",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get filter from flags")
			}

			client, err := connection.NewClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			taskQueue, wait := core.WorkerPool(ctx, 16)
			defer wait()

			// Generate tasks.
			name := args[0]
			if spec, err := names.ParseSpec(name); err == nil {
				// Iterate through a collection of specs and evaluate each.
				err = core.ListSpecs(ctx, client, spec, filter, func(spec *rpc.ApiSpec) {
					taskQueue <- &generateOpenAPITask{
						client:   client,
						specName: spec.Name,
					}
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list specs")
				}
			}
		},
	}
}

type generateOpenAPITask struct {
	client   connection.Client
	specName string
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
	data, err := core.GetBytesForSpec(ctx, task.client, spec)
	if err != nil {
		return err
	}
	relation := "generated.yaml"
	var openapi string
	if core.IsProto(spec.GetMimeType()) && core.IsZipArchive(spec.GetMimeType()) {
		log.FromContext(ctx).Debugf("Computing %s/specs/%s", spec.Name, relation)
		openapi, err = openAPIFromZippedProtos(spec.Name, data)
		if err != nil {
			log.FromContext(ctx).WithError(err).Warnf("error processing protos: %s", spec.Name)
		}
	} else {
		log.FromContext(ctx).Infof("we don't know how to generate OpenAPI for %s", spec.Name)
		return nil
	}
	specName, _ := names.ParseSpec(spec.GetName())
	messageData := []byte(openapi)
	newSpec := &rpc.ApiSpec{
		Name:     specName.Version().Spec(relation).String(),
		MimeType: core.MimeTypeForMessageType("application/x.openapi;version=3"),
		Contents: messageData,
	}
	_, err = task.client.UpdateApiSpec(ctx, &rpc.UpdateApiSpecRequest{
		ApiSpec:      newSpec,
		AllowMissing: true,
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Warn("Failed to generate OpenAPI")
	}
	return nil
}

// openAPIFromZippedProtos runs the OpenAPI generator and returns the results.
func openAPIFromZippedProtos(name string, b []byte) (string, error) {
	// create a tmp directory
	root, err := ioutil.TempDir("", "registry-protos-")
	if err != nil {
		return "", err
	}
	fmt.Printf("WORKING IN %s\n", root)
	// whenever we finish, delete the tmp directory
	defer os.RemoveAll(root)
	// unzip the protos to the temp directory
	_, err = core.UnzipArchiveToPath(b, root+"/protos")
	if err != nil {
		return "", err
	}
	// unpack api-common-protos in the temp directory
	cmd := exec.Command("git", "clone", "https://github.com/googleapis/api-common-protos")
	cmd.Dir = root
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	// run on each proto file in the archive
	lint, err := generateOpenAPIForDirectory(name, root)
	if err == nil {
		return lint, nil
	}
	// if we had errors, add googleapis to the temp directory and retry
	cmd = exec.Command("git", "clone", "https://github.com/googleapis/googleapis")
	cmd.Dir = root
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	// rerun with the extra googleapis protos
	return generateOpenAPIForDirectory(name, root)
}

func generateOpenAPIForDirectory(name string, root string) (string, error) {
	lint := &rpc.Lint{}
	lint.Name = name
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
	fmt.Printf("protoc should be run on %+v\n", protos)
	parts := []string{}
	parts = append(parts, protos...)
	parts = append(parts, "-I")
	parts = append(parts, "protos")
	parts = append(parts, "-I")
	parts = append(parts, "api-common-protos")
	parts = append(parts, "-I")
	parts = append(parts, "googleapis")
	parts = append(parts, "--openapi_out=.")
	cmd := exec.Command("protoc", parts...)
	cmd.Dir = root
	fmt.Printf("running %+v\n", cmd)
	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error %+v\n", err)
		return "", err
	}
	fmt.Printf("output: %s\n", string(data))
	// attempt to read an openapi.yaml file
	bytes, err := ioutil.ReadFile(root + "/openapi.yaml")
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
