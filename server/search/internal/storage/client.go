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

package storage

import (
	"context"

	"database/sql"

	"github.com/apigee/registry-experimental/server/search/internal/storage/gorm"
)

type Client struct {
	*gorm.Client
}

func NewClient(ctx context.Context, driver, dsn string) (*Client, error) {
	gc, err := gorm.NewClient(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: gc,
	}, nil
}

type RawRows interface {
	Append(rows *sql.Rows) error
}
