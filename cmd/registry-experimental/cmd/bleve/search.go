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

package bleve

import (
	"fmt"

	"github.com/apigee/registry/pkg/log"
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/spf13/cobra"
)

func searchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search a local search index",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			// open an existing index
			index, err := bleve.Open("registry.bleve")
			if err != nil {
				log.FromContext(ctx).WithError(err).Debug("Failed to open bleve")
				return
			}

			// search for some text
			query := bleve.NewQueryStringQuery(args[0])
			search := bleve.NewSearchRequest(query)
			search.Highlight = bleve.NewHighlightWithStyle("ansi")
			searchResults, err := index.Search(search)
			if err != nil {
				log.FromContext(ctx).WithError(err).Debug("Failed to search index")
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%+v", searchResults)
		},
	}
}
