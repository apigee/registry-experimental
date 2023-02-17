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

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover"
	"github.com/apigee/registry/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Version value will be replaced by the release tag when the binaries are
// generated by GoReleaser.
var Version = "dev"

var verbose bool

// func init() {
// 	pflag.BoolVarP(&verbose, "verbose", "v", false, "Print YAML to stdout before writing to disk")
// }

func main() {
	// Bind a logger instance to the local context with metadata for outbound requests.
	logger := log.NewLogger(log.DebugLevel)
	ctx := log.NewOutboundContext(log.NewContext(context.Background(), logger), log.Metadata{
		UID: fmt.Sprintf("%.8s", uuid.New()),
	})

	cmd := &cobra.Command{
		Use:     "registry-connect",
		Version: Version,
		Short:   "Exports Apigee resources to YAML files compatible with API Registry.",
		Long: `This command exports API Registry-compatible YAML files from an Apigee
instance into the specified DIRECTORY. (If no DIRECTORY is specified, it will print 
to the console.)
		
Once this command has been successfully run, the entire exported DIRECTORY or
individual files can be imported into an API Registry instance by running:

  registry apply -f DIRECTORY|FILE

See "registry apply --help" for more information.`,
	}

	cmd.AddCommand(discover.Command())

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println("----")
		fmt.Println("Need more help?")
		fmt.Println("https://github.com/apigee/registry/wiki")
		os.Exit(1)
	}
}
