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

package cmd

import (
	"context"

	"github.com/apigee/registry-experimental/cmd/zero/cmd/apikeys"
	"github.com/apigee/registry-experimental/cmd/zero/cmd/petstore"
	"github.com/apigee/registry-experimental/cmd/zero/cmd/servicecontrol"
	"github.com/apigee/registry-experimental/cmd/zero/cmd/servicemanagement"
	"github.com/spf13/cobra"
)

func Cmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "zero",
	}

	cmd.AddCommand(apikeys.Cmd())
	cmd.AddCommand(petstore.Cmd())
	cmd.AddCommand(servicecontrol.Cmd())
	cmd.AddCommand(servicemanagement.Cmd())
	return cmd
}
