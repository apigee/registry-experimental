// Copyright 2023 Google LLC.
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
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a simple search API using a local bleve index",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			router := gin.Default()
			router.GET("/search", search)
			router.POST("/index", index)
			return router.Run(fmt.Sprintf("0.0.0.0:%d", port))
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 8888, "port for server")
	return cmd
}

func search(c *gin.Context) {
	q := c.Query("q")

	limit := 10
	if c.Query("l") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("l"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,
				gin.H{"message": fmt.Sprintf("invalid limit: %s", err)})
			return
		}
	}
	index, err := bleve.Open(bleveDir)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"message": fmt.Sprintf("failed to open search index: %s", err)})
		return
	}
	defer index.Close()
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequest(query)
	search.Size = limit
	search.Highlight = bleve.NewHighlightWithStyle("ansi")
	searchResults, err := index.Search(search)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"message": fmt.Sprintf("failed to search index: %s", err)})
		return
	}
	c.IndentedJSON(http.StatusOK, searchResults)
}

type IndexRequestBody struct {
	Pattern string `json:"pattern"`
	Filter  string `json:"filter"`
}

func index(c *gin.Context) {
	var requestBody IndexRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("%s", err)
	} else {
		log.Printf("pattern %+v", requestBody.Pattern)
		log.Printf("filter %+v", requestBody.Filter)
		cmd := indexCommand()
		args := []string{requestBody.Pattern}
		if requestBody.Filter != "" {
			args = append(args, "--filter")
			args = append(args, requestBody.Filter)
		}
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			log.Printf("%s", err)
		}
	}
	c.IndentedJSON(http.StatusOK, []string{})
}
