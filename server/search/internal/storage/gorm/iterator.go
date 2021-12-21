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

package gorm

import (
	"encoding/base64"
	"fmt"

	"github.com/apigee/registry-experimental/server/search/internal/storage/models"
	"google.golang.org/api/iterator"
)

// Iterator can be used to iterate through results of a query.
type Iterator struct {
	Client *Client
	Values interface{}
	Index  int
	Cursor string
}

// GetCursor gets the cursor for the next page of results.
func (it *Iterator) GetCursor() (string, error) {
	encodedCursor := base64.StdEncoding.EncodeToString([]byte(it.Cursor))
	return encodedCursor, nil
}

// Next gets the next value from the iterator.
func (it *Iterator) Next(v interface{}) (*Key, error) {
	switch x := v.(type) {
	case *models.Document:
		values := it.Values.([]models.Document)
		if it.Index < len(values) {
			*x = values[it.Index]
			it.Cursor = x.Key
			it.Index++
			return it.Client.NewKey("Document", x.Key), nil
		}
		return nil, iterator.Done
	default:
		return nil, fmt.Errorf("unsupported iterator type: %t", v)
	}
}
