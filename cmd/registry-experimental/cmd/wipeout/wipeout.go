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

package wipeout

import (
	"context"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/log"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
)

func Command(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wipeout PROJECT-ID",
		Short: "Delete everything in a project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			registryClient, err := connection.NewClient(ctx)
			if err != nil {
				log.Fatalf(ctx, "Failed to create client: %+v", err)
			}
			defer registryClient.Close()
			jobs, err := cmd.Flags().GetInt("jobs")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get jobs from flags")
			}
			taskQueue, wait := core.WorkerPool(ctx, jobs)
			defer wait()

			log.Infof(ctx, "Deleting everything in project %s", args[0])
			wipeoutProject(ctx, registryClient, taskQueue, "projects/"+args[0]+"/locations/global")
		},
	}
	cmd.Flags().Int("jobs", 10, "Number of deletion requests to make simultaneously")
	return cmd
}

type ResourceKind int

const (
	API ResourceKind = iota
	Version
	Spec
	Deployment
	Artifact
)

type deleteResourceTask struct {
	client connection.Client
	name   string
	kind   ResourceKind
}

func (task *deleteResourceTask) String() string {
	return "delete " + task.name
}

func (task *deleteResourceTask) Run(ctx context.Context) error {
	var err error
	switch task.kind {
	case API:
		log.Infof(ctx, "delete API %s", task.name)
		err = task.client.DeleteApi(ctx, &rpc.DeleteApiRequest{Name: task.name})
	case Version:
		log.Infof(ctx, "delete version %s", task.name)
		err = task.client.DeleteApiVersion(ctx, &rpc.DeleteApiVersionRequest{Name: task.name})
	case Spec:
		log.Infof(ctx, "delete spec %s", task.name)
		err = task.client.DeleteApiSpec(ctx, &rpc.DeleteApiSpecRequest{Name: task.name})
	case Deployment:
		log.Infof(ctx, "delete deployment %s", task.name)
		err = task.client.DeleteApiDeployment(ctx, &rpc.DeleteApiDeploymentRequest{Name: task.name})
	case Artifact:
		log.Infof(ctx, "delete artifact %s", task.name)
		err = task.client.DeleteArtifact(ctx, &rpc.DeleteArtifactRequest{Name: task.name})
	default:
		log.Infof(ctx, "unknown resource type %s", task.name)
	}
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Deletion failed")
	}
	return nil
}

func wipeoutProject(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, project string) {
	wipeoutArtifacts(ctx, registryClient, taskQueue, project+"/apis/-/versions/-/specs/-")
	wipeoutArtifacts(ctx, registryClient, taskQueue, project+"/apis/-/versions/-")
	wipeoutArtifacts(ctx, registryClient, taskQueue, project+"/apis/-/deployments/-")
	wipeoutArtifacts(ctx, registryClient, taskQueue, project+"/apis/-")
	// don't wipeout project-level artifacts
	// wipeoutArtifacts(ctx, registryClient, taskQueue, project)
	wipeoutApiSpecs(ctx, registryClient, taskQueue, project+"/apis/-/versions/-")
	wipeoutApiVersions(ctx, registryClient, taskQueue, project+"/apis/-")
	wipeoutApiDeployments(ctx, registryClient, taskQueue, project+"/apis/-")
	wipeoutApis(ctx, registryClient, taskQueue, project)
}

func wipeoutArtifacts(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, parent string) {
	it := registryClient.ListArtifacts(ctx, &rpc.ListArtifactsRequest{Parent: parent})
	for artifact, err := it.Next(); err == nil; artifact, err = it.Next() {
		taskQueue <- &deleteResourceTask{
			client: registryClient,
			name:   artifact.Name,
			kind:   Artifact,
		}
	}
}

func wipeoutApiDeployments(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, parent string) {
	it := registryClient.ListApiDeployments(ctx, &rpc.ListApiDeploymentsRequest{Parent: parent})
	for deployment, err := it.Next(); err == nil; deployment, err = it.Next() {
		taskQueue <- &deleteResourceTask{
			client: registryClient,
			name:   deployment.Name,
			kind:   Deployment,
		}
	}
}

func wipeoutApiSpecs(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, parent string) {
	it := registryClient.ListApiSpecs(ctx, &rpc.ListApiSpecsRequest{Parent: parent})
	for spec, err := it.Next(); err == nil; spec, err = it.Next() {
		taskQueue <- &deleteResourceTask{
			client: registryClient,
			name:   spec.Name,
			kind:   Spec,
		}
	}
}

func wipeoutApiVersions(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, parent string) {
	it := registryClient.ListApiVersions(ctx, &rpc.ListApiVersionsRequest{Parent: parent})
	for version, err := it.Next(); err == nil; version, err = it.Next() {
		taskQueue <- &deleteResourceTask{
			client: registryClient,
			name:   version.Name,
			kind:   Version,
		}
	}
}

func wipeoutApis(ctx context.Context, registryClient connection.Client, taskQueue chan<- core.Task, parent string) {
	it := registryClient.ListApis(ctx, &rpc.ListApisRequest{Parent: parent})
	for api, err := it.Next(); err == nil; api, err = it.Next() {
		taskQueue <- &deleteResourceTask{
			client: registryClient,
			name:   api.Name,
			kind:   API,
		}
	}
}
