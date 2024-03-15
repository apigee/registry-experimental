// Copyright 2023 Google LLC. All Rights Reserved.
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

package petstore

import (
	"net/http"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/servicecontrol"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func serveCmd() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:  "serve SERVICE",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			router := gin.Default()
			router.Use(servicecontrol.Middleware(args[0], verbose))
			router.GET("/v1/pets", GetPets)
			router.GET("/v1/pets/:id", GetPetById)
			router.POST("/v1/pets", CreatePet)
			return router.Run("0.0.0.0:8080")
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging.")
	return cmd
}

// pet represents data about a pet.
type pet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// pets slice to seed pet data.
var pets = []pet{
	{ID: "1", Name: "Tardar Sauce", Tag: "cat"},
	{ID: "2", Name: "Bo", Tag: "dog"},
	{ID: "3", Name: "Toto", Tag: "dog"},
}

// GetPets responds with the list of all pets as JSON.
func GetPets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, pets)
}

// CreatePet adds a pet from JSON received in the request body.
func CreatePet(c *gin.Context) {
	var newpet pet

	// Call BindJSON to bind the received JSON to newpet.
	if err := c.BindJSON(&newpet); err != nil {
		return
	}

	// Add the new pet to the slice.
	pets = append(pets, newpet)
	c.IndentedJSON(http.StatusCreated, newpet)
}

// GetPetById locates the pet whose ID value matches the id
// parameter sent by the client, then returns that pet as a response.
func GetPetById(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of pets, looking for
	// a pet whose ID value matches the parameter.
	for _, a := range pets {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "pet not found",
	})
}
