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
	"regexp"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/spf13/cobra"
)

var limit int
var output string

func searchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search QUERY",
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
			search.Size = limit
			search.Highlight = bleve.NewHighlightWithStyle("ansi")
			searchResults, err := index.Search(search)
			if err != nil {
				return fmt.Errorf("failed to search index: %s", err)
			}
			switch output {
			case "json":
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
			case "text":
				fmt.Fprintf(cmd.OutOrStdout(), "%d total hits (%s)\n", searchResults.Total, searchResults.Took)
				for i, hit := range searchResults.Hits {
					fmt.Fprintf(cmd.OutOrStdout(), "\n>> %d (%0.2f) %s\n", i+1, hit.Score, reduce(hit.ID))
					for k, v := range hit.Fragments {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n%s\n", k, v)
					}
				}
			case "name":
				fmt.Fprintf(cmd.OutOrStdout(), "%d total hits (%s)\n", searchResults.Total, searchResults.Took)
				for i, hit := range searchResults.Hits {
					fmt.Fprintf(cmd.OutOrStdout(), ">> %d (%0.2f) %s\n", i+1, hit.Score, reduce(hit.ID))
				}
			default:
				return fmt.Errorf("unknown output format: %q", output)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "text", "output type (text|json|name)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "maximum number of results to return")
	return cmd
}

// reduce shortens hit ids to project-local names.
func reduce(name string) string {
	re := regexp.MustCompile(`^projects\/.*\/locations\/global\/(.*)$`)
	matches := re.FindAllStringSubmatch(name, 1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return name
}
