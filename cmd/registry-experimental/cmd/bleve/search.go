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
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/spf13/cobra"
)

func searchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search a local search index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// open an existing index
			index, err := bleve.Open(bleveDir)
			if err != nil {
				return fmt.Errorf("failed to open search index: %s", err)
			}
			defer index.Close()

			// search for some text
			query := bleve.NewQueryStringQuery(args[0])
			search := bleve.NewSearchRequest(query)
			search.Highlight = bleve.NewHighlightWithStyle("ansi")
			searchResults, err := index.Search(search)
			if err != nil {
				return fmt.Errorf("failed to search index: %s", err)
			}
			bytes, err := json.MarshalIndent(searchResults, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to serialize search results: %s", err)
			}
			if _, err = cmd.OutOrStdout().Write(bytes); err != nil {
				return err
			}
			if _, err = cmd.OutOrStdout().Write([]byte("\n")); err != nil {
				return err
			}
			return nil
		},
	}
}
