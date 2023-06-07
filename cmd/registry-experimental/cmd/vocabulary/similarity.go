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

package vocabulary

import (
	"context"
	"fmt"

	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	metrics "github.com/google/gnostic/metrics"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var filter string

func similarityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "similarity PATTERN1...",
		Short: "Compute similarities of vocabularies",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				return err
			}
			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				return err
			}
			vocabularyArtifacts := make([]*rpc.Artifact, 0)
			for _, p := range args {
				pattern := c.FQName(p)
				name, err := names.ParseArtifact(pattern)
				if err != nil {
					return err
				}
				if err = visitor.ListArtifacts(ctx, client,
					name,
					filter,
					true, func(ctx context.Context, message *rpc.Artifact) error {
						messageType, err := mime.MessageTypeForMimeType(message.GetMimeType())
						if err != nil || messageType != "gnostic.metrics.Vocabulary" {
							log.Debugf(ctx, "Skipping, not a vocabulary: %s", message.Name)
							return nil
						}
						vocabularyArtifacts = append(vocabularyArtifacts, message)
						return nil
					}); err != nil {
					return err
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), ",")
			for i := 0; i < len(vocabularyArtifacts); i++ {
				fmt.Fprintf(cmd.OutOrStdout(), "%s,", shortName(vocabularyArtifacts[i].Name))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n")
			for i := 0; i < len(vocabularyArtifacts); i++ {
				vi := &metrics.Vocabulary{}
				if err := proto.Unmarshal(vocabularyArtifacts[i].GetContents(), vi); err != nil {
					log.FromContext(ctx).WithError(err).Debug("Failed to unmarshal contents")
					return nil
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s,", shortName(vocabularyArtifacts[i].Name))
				for j := 0; j < i; j++ {
					fmt.Fprintf(cmd.OutOrStdout(), "%1.3f,", 0.0)
				}
				for j := i; j < len(vocabularyArtifacts); j++ {
					vj := &metrics.Vocabulary{}
					if err := proto.Unmarshal(vocabularyArtifacts[j].GetContents(), vj); err != nil {
						log.FromContext(ctx).WithError(err).Debug("Failed to unmarshal contents")
						return nil
					}
					counts := make(map[string]int32)
					total := 0
					for _, v := range []*metrics.Vocabulary{vi, vj} {
						for _, wordlist := range [][]*metrics.WordCount{
							v.Operations,
							v.Schemas,
							v.Parameters,
							v.Properties,
						} {
							total += len(wordlist)
							for _, pair := range wordlist {
								counts[pair.Word] += pair.Count
							}
						}
					}
					// Here we are using a map that collects the total number of unique words
					// to find the number of words that are common to two collections.
					// Why does this work?
					// 		u1 = the number of words unique to v1
					// 		u2 = the number of words unique to v2
					// 		 c = the number of words common to v1 and v2
					// The number of unique words (the size of the map) is unique = u1 + u2 + c.
					unique := len(counts)
					// The total number of words in both maps is total = u1 + u2 + 2*c.
					// The number of words in common is total - merged = c.
					common := total - unique
					// We evaluate similarity as the fraction of all words that are shared.
					// similarity = 2*c / total
					similarity := 2.0 * float32(common) / float32(total)
					fmt.Fprintf(cmd.OutOrStdout(), "%1.3f,", similarity)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s", shortName(vocabularyArtifacts[i].Name))
				fmt.Fprintf(cmd.OutOrStdout(), "\n")
			}
			return err
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "filter selected resources")
	return cmd
}

func shortName(artifactName string) string {
	name, _ := names.ParseArtifact(artifactName)
	return name.ApiID() + "/" + name.VersionID()
}
