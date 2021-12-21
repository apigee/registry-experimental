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

	experimental_rpc "github.com/apigee/registry-experimental/rpc"
	"github.com/apigee/registry-experimental/server/search/internal/storage"
	registry_rpc "github.com/apigee/registry/rpc"
	"google.golang.org/genproto/googleapis/longrunning"
)

// Config configures the Search server.
type Config struct {
	Database string
	DBConfig string
}

// SearchServer implements a Search server.
type SearchServer struct {
	database string
	dbConfig string
	registry registry_rpc.RegistryServer

	experimental_rpc.UnimplementedSearchServer
	longrunning.UnimplementedOperationsServer
}

func New(config Config, r registry_rpc.RegistryServer) *SearchServer {
	return &SearchServer{
		database: config.Database,
		dbConfig: config.DBConfig,
		registry: r,
	}
}

func (s *SearchServer) getStorageClient(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx, s.database, s.dbConfig)
}
