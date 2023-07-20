// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edge

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/api/apigee/v1"
)

// https://docs.apigee.com/api-platform/analytics/use-analytics-api-measure-api-program-performance#gettingmetricswiththemanagementapi

// https://api.enterprise.apigee.com/v1/o/{org_name}/environments/{env_name}/stats/apiproxy
const metricsPath = "environments/%s/stats/%s" // env, dimensions

// MetricsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with Metrics.
type MetricsService interface {
	Metrics(ctx context.Context, env string, dimensions []string,
		metrics []string, start time.Time, end time.Time) (*apigee.GoogleCloudApigeeV1Stats, *Response, error)
}

// MetricsServiceOp represents metrics
type MetricsServiceOp struct {
	client *EdgeClient
}

var _ MetricsService = &MetricsServiceOp{}

func (s *MetricsServiceOp) Metrics(ctx context.Context, env string, dimensions []string,
	metrics []string, start time.Time, end time.Time) (*apigee.GoogleCloudApigeeV1Stats, *Response, error) {

	ds := strings.Join(dimensions, ",")
	ms := strings.Join(metrics, ",")
	timeFormat := "01/02/2006 15:04"
	tr := start.Format(timeFormat) + "~" + end.Format(timeFormat)

	q := url.Values{}
	q.Set("select", ms)
	q.Set("timeRange", tr)
	urlString := fmt.Sprintf(metricsPath, env, ds) + "?" + q.Encode()
	req, err := s.client.NewRequestNoEnv("GET", urlString, nil)
	if err != nil {
		return nil, nil, err
	}

	stats := apigee.GoogleCloudApigeeV1Stats{}
	resp, err := s.client.Do(req, &stats)
	if err != nil {
		return nil, resp, err
	}
	return &stats, resp, err
}
