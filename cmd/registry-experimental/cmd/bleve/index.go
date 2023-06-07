// Copyright 2020 Google LLC. All Rights Reserved.
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

package bleve

import (
	"context"
	"fmt"
	"sync"

	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/cmd/registry/patch"
	"github.com/apigee/registry/cmd/registry/tasks"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/blevesearch/bleve"
	"github.com/spf13/cobra"
)

var bleveMutex sync.Mutex

var filter string
var jobs int

func indexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index PATTERN",
		Short: "Add documents to a local search index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to get client: %s", err)
			}
			c, err := connection.ActiveConfig()
			if err != nil {
				return err
			}
			pattern := c.FQName(args[0])
			// Initialize task queue.
			taskQueue, wait := tasks.WorkerPoolIgnoreError(ctx, jobs)
			defer wait()
			v := &indexVisitor{
				taskQueue:      taskQueue,
				registryClient: registryClient,
			}
			return visitor.Visit(ctx, v, visitor.VisitorOptions{
				RegistryClient: registryClient,
				Pattern:        pattern,
				Filter:         filter,
			})
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "filter selected resources")
	cmd.Flags().IntVarP(&jobs, "jobs", "j", 10, "number of actions to perform concurrently")
	return cmd
}

type indexVisitor struct {
	visitor.Unsupported
	taskQueue      chan<- tasks.Task
	registryClient connection.RegistryClient
}

func (v *indexVisitor) SpecHandler() visitor.SpecHandler {
	return func(ctx context.Context, message *rpc.ApiSpec) error {
		v.taskQueue <- &indexSpecTask{
			client: v.registryClient,
			name:   message.Name,
		}
		return nil
	}
}

func (v *indexVisitor) ArtifactHandler() visitor.ArtifactHandler {
	return func(ctx context.Context, message *rpc.Artifact) error {
		v.taskQueue <- &indexArtifactTask{
			client: v.registryClient,
			name:   message.Name,
		}
		return nil
	}
}

type indexSpecTask struct {
	client connection.RegistryClient
	name   string
}

func (task *indexSpecTask) String() string {
	return "index " + task.name
}

func (task *indexSpecTask) Run(ctx context.Context) error {
	request := &rpc.GetApiSpecRequest{
		Name: task.name,
	}
	spec, err := task.client.GetApiSpec(ctx, request)
	if err != nil {
		return err
	}
	err = visitor.FetchSpecContents(ctx, task.client, spec)
	if err != nil {
		return nil
	}
	data := spec.GetContents()
	var message interface{}
	switch {
	case spec.GetMimeType() == "text/plain" ||
		mime.IsOpenAPIv2(spec.GetMimeType()) ||
		mime.IsOpenAPIv3(spec.GetMimeType()) ||
		mime.IsDiscovery(spec.GetMimeType()):
		message = map[string]string{spec.GetFilename(): string(data)}
	case mime.IsZipArchive(spec.GetMimeType()):
		m, err := compress.UnzipArchiveToMap(data)
		if err != nil {
			return err
		}
		// for bleve, the map content must be strings
		m2 := make(map[string]string)
		for k, v := range m {
			m2[k] = string(v)
		}
		message = m2
	default:
		return fmt.Errorf("unable to index style %s", spec.GetMimeType())
	}

	// The bleve index requires serialized updates.
	bleveMutex.Lock()
	defer bleveMutex.Unlock()
	// Open the index, creating a new one if necessary.
	index, err := bleve.Open(bleveDir)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(bleveDir, mapping)
		if err != nil {
			return err
		}
	}
	defer index.Close()
	// Index the spec.
	log.Infof(ctx, "Indexing %s", task.name)
	return index.Index(task.name, message)
}

type indexArtifactTask struct {
	client connection.RegistryClient
	name   string
}

func (task *indexArtifactTask) String() string {
	return "index " + task.name
}

func (task *indexArtifactTask) Run(ctx context.Context) error {
	request := &rpc.GetArtifactRequest{
		Name: task.name,
	}
	artifact, err := task.client.GetArtifact(ctx, request)
	if err != nil {
		return err
	}
	messageType, err := mime.MessageTypeForMimeType(artifact.GetMimeType())
	if err != nil {
		return err
	}
	switch messageType {
	case "google.cloud.apigeeregistry.v1.apihub.FieldSet":
		// supported
	default:
		return fmt.Errorf("unable to index type %s", messageType)
	}
	err = visitor.FetchArtifactContents(ctx, task.client, artifact)
	if err != nil {
		return nil
	}
	message, err := mime.MessageForMimeType(artifact.GetMimeType())
	if err != nil {
		return nil
	}
	if err := patch.UnmarshalContents(artifact.GetContents(), artifact.GetMimeType(), message); err != nil {
		return nil
	}
	// The bleve index requires serialized updates.
	bleveMutex.Lock()
	defer bleveMutex.Unlock()
	// Open the index, creating a new one if necessary.
	index, err := bleve.Open(bleveDir)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(bleveDir, mapping)
		if err != nil {
			return err
		}
	}
	defer index.Close()
	// Index the spec.
	log.Infof(ctx, "Indexing %s", task.name)
	return index.Index(task.name, message)
}
