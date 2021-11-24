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
	"os/exec"

	"github.com/apex/log"
	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"github.com/spf13/cobra"
)

func clientCommand(ctx context.Context) *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Generate a client library for an API specification",
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
					taskQueue <- &generateClientTask{
						client:   client,
						specName: spec.Name,
						language: language,
					}
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list specs")
				}
			}
		},
	}

	cmd.Flags().StringVar(&language, "language", "", "The language of the generated client")
	return cmd
}

type generateClientTask struct {
	client   connection.Client
	specName string
	language string
}

func (task *generateClientTask) String() string {
	return fmt.Sprintf("generate %s client for %s", task.language, task.specName)
}

func (task *generateClientTask) Run(ctx context.Context) error {
	log.FromContext(ctx).Info(task.String())
	request := &rpc.GetApiSpecRequest{
		Name: task.specName,
	}
	spec, err := task.client.GetApiSpec(ctx, request)
	if err != nil {
		return err
	}
	var client []byte
	relation := "client"
	if core.IsDiscovery(spec.GetMimeType()) {
		log.FromContext(ctx).Debugf("Computing %s/specs/%s", spec.Name, relation)
		data, err := core.GetBytesForSpec(ctx, task.client, spec)
		if err != nil {
			return err
		}
		client, err = clientFromDiscovery(spec.Name, data, task.language)
		if err != nil {
			log.FromContext(ctx).WithError(err).Warnf("error generating client for: %s", spec.Name)
			return nil
		}
	} else {
		log.FromContext(ctx).Infof("we don't know how to generate OpenAPI for %s", spec.Name)
		return nil
	}
	subject := spec.GetName()
	artifact := &rpc.Artifact{
		Name:     subject + "/artifacts/" + relation,
		MimeType: "application/octet-stream",
		Contents: client,
	}
	err = core.SetArtifact(ctx, task.client, artifact)
	if err != nil {
		return err
	}
	return nil
}

// clientFromDiscovery runs a discovery code generator and returns the results.
func clientFromDiscovery(name string, b []byte, language string) ([]byte, error) {
	// create a tmp directory
	root, err := ioutil.TempDir("", "registry-generator-")
	if err != nil {
		return nil, err
	}
	fmt.Printf("WORKING IN %s\n", root)
	// whenever we finish, delete the tmp directory
	//defer os.RemoveAll(root)

	// write the spec into the temp directory
	ioutil.WriteFile(root+"/discovery.json", b, 0444)

	// run the code generator
	cmd := exec.Command("google-api-go-generator",
		"-api_json_file",
		"discovery.json",
		"-gendir",
		"client")
	cmd.Dir = root
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// zip the result
	cmd = exec.Command("zip", "-r", "client.zip", "client")
	cmd.Dir = root
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(root + "/client.zip")
}
