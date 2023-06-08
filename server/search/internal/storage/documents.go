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

	"github.com/apigee/registry-experimental/server/search/internal/storage/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DocumentRows struct {
	Rows []models.Document
}

func (s *DocumentRows) Append(rows *sql.Rows) error {
	for rows.Next() {
		var row models.Document
		if err := rows.Scan(&row.Key, &row.Raw); err != nil {
			return err
		}
		s.Rows = append(s.Rows, row)
	}
	return nil
}

func (d *Client) UpdateDocuments(ctx context.Context, documents []*models.Document) error {
	for _, x := range documents {
		if x.IsEmpty() {
			if err := d.DeleteDocument(ctx, x); err != nil {
				return err
			}
		} else {
			if err := d.SaveDocument(ctx, x); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Client) DeleteDocument(ctx context.Context, document *models.Document) error {
	q := d.NewQuery(models.DocumentEntityName)
	q = q.Require("Key", document.Key)
	if err := d.Delete(ctx, q); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (d *Client) SaveDocument(ctx context.Context, document *models.Document) error {
	if err := d.PutDocument(ctx, document); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

const textSearchQuery = `
SELECT
  key, ts_headline(raw, q, 'StartSel="**", StopSel="**"') as excerpt
FROM (
  SELECT
    key, raw, ts_rank(vector, q) as rank, q
  FROM
    documents, plainto_tsquery(?) q
  WHERE
    vector @@ q
  ORDER BY
    rank DESC
) AS subquery_because_ts_headline_is_expensive_and_subquery_needs_a_name
`

func (d *Client) ListDocuments(ctx context.Context, query string) (*DocumentRows, error) {
	var rows DocumentRows
	if err := d.Raw(ctx, &rows, textSearchQuery, query); err != nil {
		return nil, err
	}
	return &rows, nil
}
