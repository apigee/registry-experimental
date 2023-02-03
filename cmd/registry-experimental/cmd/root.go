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

package cmd

import (
	"context"

	"fmt"

	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/compute"
	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/count"
	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/export"
	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/generate"
	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/search"
	"github.com/apigee/registry-experimental/cmd/registry-experimental/cmd/wipeout"
	"github.com/apigee/registry/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func Command(ctx context.Context) *cobra.Command {
	var logID string
	var cmd = &cobra.Command{
		Use:   "registry-experimental",
		Short: "Experimental utilities for working with the API Registry",
	}

	// Bind a logger instance to the local context with metadata for outbound requests.
	logger := log.NewLogger(log.DebugLevel)
	ctx = log.NewOutboundContext(log.NewContext(ctx, logger), log.Metadata{
		UID: fmt.Sprintf("%.8s", uuid.New()),
	})

	cmd.AddCommand(compute.Command(ctx))
	cmd.AddCommand(generate.Command(ctx))
	cmd.AddCommand(count.Command())
	cmd.AddCommand(export.Command())
	cmd.AddCommand(search.Command(ctx))
	cmd.AddCommand(wipeout.Command(ctx))

	cmd.PersistentFlags().StringVar(&logID, "log-id", "", "Assign an ID which gets attached to the log produced")
	return cmd
}
