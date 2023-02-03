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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/cmd/registry/types"
	"github.com/apigee/registry/log"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	discovery "github.com/google/gnostic/discovery"
	oas2 "github.com/google/gnostic/openapiv2"
	oas3 "github.com/google/gnostic/openapiv3"
	"google.golang.org/protobuf/types/descriptorpb"
)

func descriptorCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "descriptor",
		Short: "Compute descriptors of API specs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get filter from flags")
			}

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			taskQueue, wait := core.WorkerPool(ctx, 1)
			defer wait()
			// Generate tasks.
			name := args[0]
			if spec, err := names.ParseSpec(name); err == nil {
				err = visitor.ListSpecs(ctx, client, spec, filter, false, func(spec *rpc.ApiSpec) error {
					taskQueue <- &computeDescriptorTask{
						client:   client,
						specName: spec.Name,
					}
					return nil
				})
				if err != nil {
					log.FromContext(ctx).WithError(err).Fatal("Failed to list specs")
				}
			}
		},
	}
}

type computeDescriptorTask struct {
	client   connection.RegistryClient
	specName string
}

func (task *computeDescriptorTask) String() string {
	return "compute descriptor " + task.specName
}

func (task *computeDescriptorTask) Run(ctx context.Context) error {
	request := &rpc.GetApiSpecRequest{
		Name: task.specName,
	}
	spec, err := task.client.GetApiSpec(ctx, request)
	if err != nil {
		return err
	}
	name := spec.GetName()
	relation := "descriptor"
	log.Infof(ctx, "Computing %s/artifacts/%s", name, relation)
	data, err := core.GetBytesForSpec(ctx, task.client, spec)
	if err != nil {
		return nil
	}
	subject := spec.GetName()
	var typeURL string
	var document proto.Message
	if types.IsOpenAPIv2(spec.GetMimeType()) {
		typeURL = "gnostic.openapiv2.Document"
		document, err = oas2.ParseDocument(data)
		if err != nil {
			return err
		}
	} else if types.IsOpenAPIv3(spec.GetMimeType()) {
		typeURL = "gnostic.openapiv3.Document"
		document, err = oas3.ParseDocument(data)
		if err != nil {
			return err
		}
	} else if types.IsDiscovery(spec.GetMimeType()) {
		typeURL = "gnostic.discoveryv1.Document"
		document, err = discovery.ParseDocument(data)
		if err != nil {
			return err
		}
	} else if types.IsProto(spec.GetMimeType()) && types.IsZipArchive(spec.GetMimeType()) {
		typeURL = "google.protobuf.FileDescriptorSet"
		document, err = descriptorFromZippedProtos(ctx, spec.Name, data)
		if err != nil {
			log.FromContext(ctx).WithError(err).Warnf("error processing protos: %s", spec.Name)
		}
	} else {
		return fmt.Errorf("unable to generate descriptor for style %s", spec.GetMimeType())
	}
	messageData, err := proto.Marshal(document)
	if err != nil {
		return err
	}
	// TODO: consider gzipping descriptors to reduce size;
	// this will probably require some representation of compression type in the typeURL
	artifact := &rpc.Artifact{
		Name:     subject + "/artifacts/" + relation,
		MimeType: types.MimeTypeForMessageType(typeURL),
		Contents: messageData,
	}
	return core.SetArtifact(ctx, task.client, artifact)
}

// descriptorFromZippedProtos runs protoc on a collection of protos and returns a file descriptor set.
func descriptorFromZippedProtos(ctx context.Context, name string, b []byte) (*descriptorpb.FileDescriptorSet, error) {
	root, err := ioutil.TempDir("", "registry-protos-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(root)
	_, err = core.UnzipArchiveToPath(b, root+"/protos")
	if err != nil {
		return nil, err
	}
	return generateDescriptorForDirectory(ctx, name, root)
}

func generateDescriptorForDirectory(ctx context.Context, name string, root string) (*descriptorpb.FileDescriptorSet, error) {
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
		return nil, err
	}
	args := []string{}
	args = append(args, protos...)
	args = append(args, "--proto_path=protos")
	args = append(args, "--descriptor_set_out=proto.pb")
	cmd := exec.Command("protoc", args...)
	cmd.Dir = root
	log.FromContext(ctx).Debugf("Running %+v", cmd)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Debugf("Output: %s", string(data))
	// attempt to read the compiler output
	bytes, err := ioutil.ReadFile(root + "/proto.pb")
	if err != nil {
		return nil, err
	}
	var s descriptorpb.FileDescriptorSet
	err = proto.Unmarshal(bytes, &s)
	return &s, err
}
