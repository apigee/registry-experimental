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
	"gopkg.in/yaml.v3"
)

func similarityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "similarity VOCABULARY1 VOCABULARY2",
		Short: "Compute similarity of two vocabularies",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				return err
			}
			name1 := c.FQName(args[0])
			name2 := c.FQName(args[1])

			client, err := connection.NewRegistryClient(ctx)
			if err != nil {
				return err
			}
			artifactName1, err := names.ParseArtifact(name1)
			if err != nil {
				return err
			}
			artifactName2, err := names.ParseArtifact(name2)
			if err != nil {
				return err
			}
			counts := make(map[string]int32)
			total := 0

			for _, artifactName := range []names.Artifact{
				artifactName1, artifactName2,
			} {
				err = visitor.GetArtifact(ctx, client, artifactName, true,
					func(ctx context.Context, artifact *rpc.Artifact) error {
						messageType, err := mime.MessageTypeForMimeType(artifact.GetMimeType())
						if err != nil || messageType != "gnostic.metrics.Vocabulary" {
							log.Debugf(ctx, "Skipping, not a vocabulary: %s", artifact.Name)
							return nil
						}
						vocab := &metrics.Vocabulary{}
						if err := proto.Unmarshal(artifact.GetContents(), vocab); err != nil {
							log.FromContext(ctx).WithError(err).Debug("Failed to unmarshal contents")
							return nil
						}
						for _, wordlist := range [][]*metrics.WordCount{
							vocab.Operations,
							vocab.Schemas,
							vocab.Parameters,
							vocab.Properties,
						} {
							total += len(wordlist)
							for _, pair := range wordlist {
								counts[pair.Word] += pair.Count
							}
						}
						return nil
					})
				if err != nil {
					return err
				}
			}

			type Metric struct {
				Total      int     `yaml:"total"`
				Merged     int     `yaml:"merged"`
				Common     int     `yaml:"common"`
				Similarity float32 `yaml:"similarity"`
			}
			// a = words unique to v1
			// b = words unique to v2
			// c = words common to v1 and v2
			// total = a + b + 2 * c
			// merged = a + b + c
			// total - merged = c
			// similarity = 2*c / total
			var m Metric
			m.Total = total
			m.Merged = len(counts)
			m.Common = m.Total - m.Merged
			m.Similarity = 2.0 * float32(m.Common) / float32(m.Total)
			b, err := yaml.Marshal(m)
			if err != nil {
				return err
			}
			var concise = true
			if concise {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "%f", m.Similarity)
			} else {
				_, err = cmd.OutOrStdout().Write(b)
			}
			return err
		},
	}
	return cmd
}
