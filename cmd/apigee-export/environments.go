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

	"google.golang.org/api/apigee/v1"
)

func envgroups(ctx context.Context, org string) ([]*apigee.GoogleCloudApigeeV1EnvironmentGroup, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Envgroups.List(org).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.EnvironmentGroups, nil
}

func attachments(ctx context.Context, group string) ([]*apigee.GoogleCloudApigeeV1EnvironmentGroupAttachment, error) {
	apg, err := apigee.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := apg.Organizations.Envgroups.Attachments.List(group).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.EnvironmentGroupAttachments, nil
}

type envMap struct {
	hostnames map[string][]string
	envgroup  map[string]string
}

func (m *envMap) Hostnames(env string) ([]string, bool) {
	if m.hostnames == nil {
		return nil, false
	}

	v, ok := m.hostnames[env]
	return v, ok
}

func (m *envMap) Envgroup(hostname string) (string, bool) {
	if m.envgroup == nil {
		return "", false
	}

	v, ok := m.envgroup[hostname]
	return v, ok
}

func newEnvMap(ctx context.Context, org string) (*envMap, error) {
	groups, err := envgroups(ctx, org)
	if err != nil {
		return nil, err
	}

	m := &envMap{
		hostnames: make(map[string][]string),
		envgroup:  make(map[string]string),
	}

	for _, group := range groups {
		envgroup := fmt.Sprintf("%s/envgroups/%s", org, group.Name)
		attachments, err := attachments(ctx, envgroup)
		if err != nil {
			return nil, err
		}

		for _, attachment := range attachments {
			for _, hostname := range group.Hostnames {
				m.hostnames[attachment.Environment] = append(m.hostnames[attachment.Environment], hostname)
				m.envgroup[hostname] = envgroup
			}
		}
	}

	return m, nil
}
