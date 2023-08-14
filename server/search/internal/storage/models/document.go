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
	"html"
)

// DocumentEntityName is used to represent documents in storage.
const DocumentEntityName = "Document"

type Field string

const (
	FieldDisplayName Field = "displayname"
	FieldDescription Field = "description"
	FieldParameters  Field = "parameters"
	FieldMethods     Field = "methods"
	FieldSchemas     Field = "schemas"
)

type Weight string

const (
	WeightA Weight = "A"
	WeightB Weight = "B"
	WeightC Weight = "C"
	WeightD Weight = "D"
)

// A document is an item stored in the search index that is a potential
// search result. A documents is typically a resource or part of a resource
// and when so, has a resource path and a possible path into the resource
// to the resource fragment corresponding to the document. Such a fragment
// might be an API operation or schema.
//
// Documents have associated text that is used to index them. This text
// is derived from properties of the resource and is stored in two forms:
// a "vector" version that is created with the Postgres ts_vector function,
// and raw text that can be presented with the results of search queries.
// This raw text can be highlighted to show search terms.
//
// The key of the document is derived from the resource name and the path
// to the resource fragment, when appropriate.
type Document struct {
	Key       string   `gorm:"primaryKey"`
	Name      string   // The name of the resource.
	Fragment  string   // The path to the fragment of the resource, possibly empty.
	Kind      string   // The type of the resource.
	Field     Field    // The field of the resource or type of the fragment, if appropriate.
	ProjectID string   // The project associated with the document and resource.
	Vector    TSVector // A Text Search Vector of the indexed text.
	Raw       string   // The raw indexed text for excerpting; has excerpt from search result.
	Escaped   bool     // If true, the raw text is escaped.
}

// Escape should be called after filling struct to HTML-escape the raw text
// in the Vector, and also then copies it to Raw.
func (x *Document) Escape() *Document {
	if x == nil || x.Escaped {
		return x
	}
	x.Vector.RawText = html.EscapeString(x.Vector.RawText)
	x.Raw = x.Vector.RawText
	x.Escaped = true
	return x
}

// IsEmpty determines whether an update should write or delete the Document.
func (x *Document) IsEmpty() bool {
	return x == nil || x.Vector.RawText == ""
}
