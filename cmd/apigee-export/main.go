// Copyright 2022 Google LLC. All Rights Reserved.
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

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var verbose bool

func init() {
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Print YAML to stdout before writing to disk")
}

func main() {
	cmd := &cobra.Command{
		Use:   "apigee-export",
		Short: "Exports Apigee resources to YAML files compatible with API Registry",
	}

	cmd.AddCommand(exportApisCommand)
	cmd.AddCommand(exportDeploymentsCommand)

	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clean(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(s)
}
