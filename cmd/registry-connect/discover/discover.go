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

package discover

import (
	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigateway"
	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Exports extenal resources to YAML files compatible with API Registry.",
		Long: `These commands export API Registry-compatible YAML files from external
sources into the specified DIRECTORY. (If no DIRECTORY is specified, it will print 
to the console.)
		
Once this command has been successfully run, the entire exported DIRECTORY or
individual files can be imported into an API Registry instance by running:

  registry apply -f DIRECTORY|FILE

See "registry apply --help" for more information.`,
		Args: cobra.NoArgs,
	}
	cmd.AddCommand(apigee.Command())
	cmd.AddCommand(apigateway.Command())
	return cmd
}
