// Copyright 2021 Google LLC.
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

	experimental_rpc "github.com/apigee/registry-experimental/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Query handles the corresponding API request.
func (s *SearchServer) Query(ctx context.Context, req *experimental_rpc.QueryRequest) (*experimental_rpc.QueryResponse, error) {
	db, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	defer db.Close()

	rows, err := db.ListDocuments(ctx, req.GetQ())
	if err != nil {
		return nil, err
	}

	var results []*experimental_rpc.QueryResponse_Result
	for _, row := range rows.Rows {
		results = append(results, &experimental_rpc.QueryResponse_Result{
			Key:     row.Key,
			Excerpt: row.Raw,
		})
	}
	return &experimental_rpc.QueryResponse{
		Results:       results,
		NextPageToken: "",
	}, nil
}
