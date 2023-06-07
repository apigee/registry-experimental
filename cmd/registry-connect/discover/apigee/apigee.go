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

package apigee

import (
	"fmt"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/products"
	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/proxies"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "apigee",
		Short: "Exports Apigee resources to YAML files compatible with API Registry.",
	}

	cmd.PersistentFlags().BoolVar(&client.Config.Debug, "debug", false, "debug mode")
	cmd.PersistentFlags().BoolVar(&client.Config.SkipVerify, "skipverify", false, "skip server certificate verify")
	cmd.PersistentFlags().StringVarP(&client.Config.Username, "username", "u", "", "username")
	cmd.PersistentFlags().StringVarP(&client.Config.Password, "password", "p", "", "password")
	cmd.PersistentFlags().StringVarP(&client.Config.MFAToken, "mfa", "", "", "multi-factor authorization token")
	cmd.PersistentFlags().BoolVarP(&client.Config.OPDK, "opdk", "", false, "edge private cloud installation")
	cmd.PersistentFlags().BoolVarP(&client.Config.Edge, "edge", "", false, "hosted edge installation")
	cmd.PersistentFlags().StringVarP(&client.Config.MgmtURL, "mgmt", "", "", "management server endpoint (opdk only)")
	cmd.PersistentFlags().StringVarP(&client.Config.Token, "token", "t", "", "Apigee OAuth or SAML token (overrides any other given credentials)")

	cmd.MarkFlagsMutuallyExclusive("opdk", "edge")
	cmd.MarkFlagsRequiredTogether("opdk", "mgmt")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if client.Config.OPDK && client.Config.MgmtURL == "" {
			return fmt.Errorf("--mgmt must be set with --opdk")
		}
		if !client.Config.OPDK && client.Config.MgmtURL != "" {
			return fmt.Errorf("--mgmt can only be set with --opdk")
		}

		return nil
	}

	cmd.AddCommand(proxies.Command())
	cmd.AddCommand(products.Command())
	return cmd
}
