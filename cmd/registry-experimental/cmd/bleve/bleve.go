// Copyright 2020 Google LLC.
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
	"github.com/spf13/cobra"
)

var bleveDir string

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bleve",
		Short: "Experimental search features built using blevesearch.com",
	}
	cmd.AddCommand(indexCommand())
	cmd.AddCommand(searchCommand())
	cmd.AddCommand(serveCommand())
	cmd.PersistentFlags().StringVar(&bleveDir, "bleve", "registry.bleve", "path to local bleve search index")
	return cmd
}
