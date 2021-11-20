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

package search

import (
	"context"

	registry_rpc "github.com/apigee/registry/rpc"

	experimental_rpc "github.com/apigee/registry-experimental/rpc"
	"github.com/apigee/registry-experimental/server/search/internal/indexer"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// Index handles the corresponding API request.
func (s *SearchServer) Index(ctx context.Context, req *experimental_rpc.IndexRequest) (*longrunning.Operation, error) {
	db, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	defer db.Close()
	spec, err := s.registry.GetApiSpec(ctx, &registry_rpc.GetApiSpecRequest{
		Name: req.ResourceName,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	contents, err := s.registry.GetApiSpecContents(ctx, &registry_rpc.GetApiSpecContentsRequest{
		Name: req.ResourceName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	documents, err := indexer.NewDocumentsForSpec(spec, contents.Data)
	if err != nil {
		return nil, err
	}
	err = db.UpdateDocuments(ctx, documents)
	if err != nil {
		return nil, err
	}
	metadata, _ := anypb.New(&experimental_rpc.IndexMetadata{})
	response, _ := anypb.New(&experimental_rpc.IndexResponse{
		Message: "OK",
	})
	return &longrunning.Operation{
		Name:     "index",
		Metadata: metadata,
		Done:     true,
		Result:   &longrunning.Operation_Response{Response: response},
	}, nil
}
