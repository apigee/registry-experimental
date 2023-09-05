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
	"sort"
	"strings"

	"github.com/apigee/registry/cmd/registry/patch"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/fatih/camelcase"
	metrics "github.com/google/gnostic/metrics"
	"github.com/kljensen/snowball"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func stemsCommand() *cobra.Command {
	var outputID string
	cmd := &cobra.Command{
		Use:   "stems",
		Short: "Compute stems for words in specified vocabularies",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get config")
			}
			pattern := c.FQName(args[0])
			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get filter from flags")
			}

			if strings.Contains(outputID, "/") {
				log.Fatal(ctx, "output-id must specify an artifact id (final segment only) and not a full name.")
			}

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}

			patternName, err := names.ParseArtifact(pattern)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Invalid pattern")
			}

			err = visitor.ListArtifacts(ctx, client, patternName, 0, filter, true, func(ctx context.Context, artifact *rpc.Artifact) error {
				messageType, err := mime.MessageTypeForMimeType(artifact.GetMimeType())
				if err != nil || messageType != "gnostic.metrics.Vocabulary" {
					log.Debugf(ctx, "Skipping, not a vocabulary: %s", artifact.Name)
					return nil
				}

				vocab := &metrics.Vocabulary{}
				if err := patch.UnmarshalContents(artifact.GetContents(), artifact.GetMimeType(), vocab); err != nil {
					log.FromContext(ctx).WithError(err).Debug("Failed to unmarshal contents")
					return nil
				}

				counts := make(map[string]int32)
				for _, wordlist := range [][]*metrics.WordCount{
					vocab.Operations,
					vocab.Schemas,
					vocab.Parameters,
					vocab.Properties,
				} {
					for _, pair := range wordlist {
						words := camelcase.Split(pair.Word)
						for _, w := range words {
							w = strings.ToLower(w)
							w, _ = snowball.Stem(w, "english", true)
							// todo: exclude numbers and special characters
							counts[w] += pair.Count
						}
					}
				}
				stems := make([]*metrics.WordCount, 0)
				for k, v := range counts {
					stems = append(stems, &metrics.WordCount{
						Word:  k,
						Count: v,
					})
				}
				// sort in decreasing order
				sort.Slice(stems, func(i, j int) bool {
					return stems[i].Count > stems[j].Count
				})
				stemVocabulary := &metrics.Vocabulary{
					Properties: stems,
				}
				if outputID != "" {
					artifactName, _ := names.ParseArtifact(artifact.Name)
					outputName := artifactName.Parent() + "/artifacts/" + outputID
					setVocabularyToArtifact(ctx, client, stemVocabulary, outputName)
				} else {
					fmt.Println(protojson.Format((vocab)))
				}
				return nil
			})
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to list artifacts")
			}
		},
	}

	cmd.Flags().StringVar(&outputID, "output-id", "stems", "artifact ID to use when saving each result vocabulary")
	return cmd
}
