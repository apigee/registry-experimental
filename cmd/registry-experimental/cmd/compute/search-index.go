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

package compute

import (
	"context"
	"fmt"
	"sync"

	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/cmd/registry/tasks"
	"github.com/apigee/registry/cmd/registry/types"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/blevesearch/bleve"
	"github.com/spf13/cobra"
)

var bleveMutex sync.Mutex

func searchIndexCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search-index",
		Short: "Compute a local search index of specs (experimental)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get filter from flags")
			}

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			taskQueue, wait := tasks.WorkerPoolWithWarnings(ctx, 64)
			defer wait()
			// Generate tasks.
			name := args[0]
			if spec, err := names.ParseSpec(name); err == nil {
				err = visitor.ListSpecs(ctx, client, spec, filter, false, func(ctx context.Context, spec *rpc.ApiSpec) error {
					taskQueue <- &indexSpecTask{
						client:   client,
						specName: spec.Name,
					}
					return nil
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list specs")
				}
			} else {
				log.FromContext(ctx).Fatalf("We don't know how to index %s", name)
			}
		},
	}
}

type indexSpecTask struct {
	client   connection.RegistryClient
	specName string
}

func (task *indexSpecTask) String() string {
	return "index " + task.specName
}

func (task *indexSpecTask) Run(ctx context.Context) error {
	request := &rpc.GetApiSpecRequest{
		Name: task.specName,
	}
	spec, err := task.client.GetApiSpec(ctx, request)
	if err != nil {
		return err
	}
	data, err := visitor.GetBytesForSpec(ctx, task.client, spec)
	if err != nil {
		return nil
	}
	var message interface{}
	switch {
	case spec.GetMimeType() == "text/plain" ||
		types.IsOpenAPIv2(spec.GetMimeType()) ||
		types.IsOpenAPIv3(spec.GetMimeType()) ||
		types.IsDiscovery(spec.GetMimeType()):
		message = map[string]string{spec.GetFilename(): string(data)}
	case types.IsProto(spec.GetMimeType()):
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
		return fmt.Errorf("unable to generate descriptor for style %s", spec.GetMimeType())
	}

	// The bleve index requires serialized updates.
	bleveMutex.Lock()
	defer bleveMutex.Unlock()
	// Open the index, creating a new one if necessary.
	const bleveDir = "registry.bleve"
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
	log.Debugf(ctx, "Indexing %s", task.specName)
	return index.Index(task.specName, message)
}
