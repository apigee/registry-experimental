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
	"regexp"
	"strings"

	"github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	apigee "github.com/apigee/registry-experimental/cmd/registry-connect/discover/apigee/client"
	"github.com/spf13/cobra"
)

var project string // TODO: remove when a relative ReferenceList_Reference.Resource works in Hub

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "apigee",
		Short: "Exports Apigee resources to YAML files compatible with API Registry.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if client.Config.OPDK && client.Config.MgmtURL == "" {
				return fmt.Errorf("--mgmt must be set with --opdk")
			}
			if !client.Config.OPDK && client.Config.MgmtURL != "" {
				return fmt.Errorf("--mgmt can only be set with --opdk")
			}

			ctx := cmd.Context()
			apigee.Config.Org = args[0]
			client, err := apigee.NewClient()
			if err != nil {
				return err
			}
			return export(ctx, client)
		},
	}
	cmd.Flags().StringVarP(&project, "project", "", "", "hub project id (temporary)")
	_ = cmd.MarkFlagRequired("project")

	cmd.Flags().BoolVar(&client.Config.Debug, "debug", false, "debug mode")
	cmd.Flags().BoolVar(&client.Config.SkipVerify, "skipverify", false, "skip server certificate verify")
	cmd.Flags().StringVarP(&client.Config.Username, "username", "u", "", "username")
	cmd.Flags().StringVarP(&client.Config.Password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&client.Config.MFAToken, "mfa", "", "", "multi-factor authorization token")
	cmd.Flags().BoolVarP(&client.Config.OPDK, "opdk", "", false, "edge private cloud installation")
	cmd.Flags().BoolVarP(&client.Config.Edge, "edge", "", false, "hosted edge installation")
	cmd.Flags().StringVarP(&client.Config.MgmtURL, "mgmt", "", "", "management server endpoint (opdk only)")
	cmd.Flags().StringVarP(&client.Config.Token, "token", "t", "", "Apigee OAuth or SAML token (overrides any other given credentials)")

	cmd.MarkFlagsMutuallyExclusive("opdk", "edge")
	cmd.MarkFlagsRequiredTogether("opdk", "mgmt")

	return cmd
}

func label(s string) string {
	return strings.ToLower(regexp.MustCompile(`([^A-Za-z0-9-_]+)`).ReplaceAllString(s, "-"))
}

func name(s string) string {
	return strings.ToLower(regexp.MustCompile(`([^A-Za-z0-9-]+)`).ReplaceAllString(s, "-"))
}
