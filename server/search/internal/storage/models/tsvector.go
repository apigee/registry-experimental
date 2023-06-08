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

package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TSVector opaquely represents the write-only ts_vector, containing
// the should-be-escaped text to index and the search weight.
type TSVector struct {
	RawText string
	Weight  Weight
}

// GormDataType of TSVector is the Postgres column type "tsvector",
func (t TSVector) GormDataType() string {
	return "tsvector"
}

// GormValue of TSVector returns the Postgres expression to convert
// search text to a weighted vector.
func (t TSVector) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	w := t.Weight
	if w == "" {
		w = WeightD
	}
	return clause.Expr{
		SQL:  "setweight(to_tsvector(?), ?)",
		Vars: []interface{}{t.RawText, w},
	}
}

// Scan implements the sql.Scanner interface
func (t *TSVector) Scan(v interface{}) error {
	t.RawText = fmt.Sprintf("%+v", v)
	return nil
}
