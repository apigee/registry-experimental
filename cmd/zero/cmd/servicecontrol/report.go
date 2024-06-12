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

package servicecontrol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"google.golang.org/api/option"
	"google.golang.org/api/servicecontrol/v1"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type distOptions struct {
	Buckets int64
	Growth  float64
	Scale   float64
}

func createInt64MetricSet(name string, value int64) *servicecontrol.MetricValueSet {
	return &servicecontrol.MetricValueSet{
		MetricName: name,
		MetricValues: []*servicecontrol.MetricValue{
			{
				Int64Value: &value,
			},
		},
	}
}

var (
	timeDistOptions = distOptions{29, 2.0, 1e-6}
	sizeDistOptions = distOptions{8, 10.0, 1}
)

func createDistMetricSet(options *distOptions, name string, value int64) *servicecontrol.MetricValueSet {
	buckets := make([]int64, options.Buckets+2)
	fValue := float64(value)
	idx := 0
	if fValue >= options.Scale {
		idx = 1 + int(math.Log(fValue/options.Scale)/math.Log(options.Growth))
		if idx >= len(buckets) {
			idx = len(buckets) - 1
		}
	}
	buckets[idx] = 1
	distValue := servicecontrol.Distribution{
		Count:        1,
		BucketCounts: buckets,
		ExponentialBuckets: &servicecontrol.ExponentialBuckets{
			NumFiniteBuckets: options.Buckets,
			GrowthFactor:     options.Growth,
			Scale:            options.Scale,
		},
	}
	if value != 0 {
		distValue.Mean = fValue
		distValue.Minimum = fValue
		distValue.Maximum = fValue
	}
	return &servicecontrol.MetricValueSet{
		MetricName: name,
		MetricValues: []*servicecontrol.MetricValue{
			{
				DistributionValue: &distValue,
			},
		},
	}
}

func reportCmd() *cobra.Command {
	var output string
	var iterations int
	cmd := &cobra.Command{
		Use:  "report",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := config.GetClient(ctx)
			if err != nil {
				return err
			}
			srv, err := servicecontrol.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				return err
			}

			for i := 0; i < iterations; i++ {
				now := time.Now()
				timestamp := now.Format(time.RFC3339Nano)
				timestampfloat := float64(now.UnixNano()) / 1e9
				time.Sleep(100 * time.Millisecond)
				timestamp2 := now.Format(time.RFC3339Nano)

				callerIP := "172.125.77.209"
				apiName := "1." + strings.ReplaceAll(
					strings.ReplaceAll(serviceName, ".", "_"),
					"-", "_")
				operationName := apiName + ".Unknown"
				uid := uuid.New().String()

				status := 200
				if i%5 == 4 {
					status = 404
				} else if i%7 == 6 {
					status = 500
				}
				operation := &servicecontrol.Operation{
					OperationId: uid,
					// operation name seems unused but cannot be empty
					OperationName: operationName,
					//ConsumerId:    "project:" + producerProject,
					ConsumerId: "api_key:" + apiKey,
					StartTime:  timestamp,
					Labels: map[string]string{
						"cloud.googleapis.com/location":             "us-west1",
						"serviceruntime.googleapis.com/api_method":  operationName,
						"cloud.googleapis.com/project":              producerProject,
						"cloud.googleapis.com/service":              serviceName,
						"serviceruntime.googleapis.com/api_version": "1.0.0",
						// none of the following appear in the logs
						// but they seem to be set in ESPv2
						//  https://github.com/GoogleCloudPlatform/esp-v2/blob/9217d68484321aceb7f7fbdc63be9363c96ed722/tests/utils/service_control_utils.go#L245
						"cloud.googleapis.com/uid":                       uid,
						"servicecontrol.googleapis.com/caller_ip":        callerIP,
						"servicecontrol.googleapis.com/service_agent":    "ESPv2/2.45.0",
						"servicecontrol.googleapis.com/platform":         "Cloud Run",
						"servicecontrol.googleapis.com/user_agent":       "ESPv2",
						"serviceruntime.googleapis.com/consumer_project": consumerProject,
						"/response_code":                                 fmt.Sprintf("%d", status),
						"/response_code_class":                           fmt.Sprintf("%dxx", status/100),
						"/status_code":                                   fmt.Sprintf("%d", status),
						"/protocol":                                      "http",
					},
				}
				{
					request := &servicecontrol.CheckRequest{
						Operation: operation,
					}
					start := time.Now()
					result, err := srv.Services.Check(serviceName, request).Do()
					elapsed := time.Since(start)
					if err != nil {
						return err
					}
					bytes, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						return &json.UnsupportedValueError{}
					}
					fmt.Printf("CHECK (%fms) %d > %s\n", float64(elapsed)/1e6, i, string(bytes))
				}

				// allocate quota
				{
					request := &servicecontrol.AllocateQuotaRequest{
						ServiceConfigId: serviceConfig,
						AllocateOperation: &servicecontrol.QuotaOperation{
							ConsumerId:  "api_key:" + apiKey,
							OperationId: uid,
							QuotaMode:   "NORMAL",
							MethodName:  "1.example5_apiregistry_dev.Unknown",
							QuotaMetrics: []*servicecontrol.MetricValueSet{
								createInt64MetricSet("calls", 1),
							},
						},
					}
					start := time.Now()
					result, err := srv.Services.AllocateQuota(serviceName, request).Do()
					elapsed := time.Since(start)
					if err != nil {
						return err
					}
					bytes, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						return &json.UnsupportedValueError{}
					}
					fmt.Printf("ALLOCATE QUOTA (%fms) %d > %s\n", float64(elapsed)/1e6, i, string(bytes))
					if result.AllocateErrors != nil {
						return errors.New("out of quota")
					}
				}

				payload := map[string]interface{}{
					"api_key_state":        "NOT CHECKED",
					"api_key":              apiKey,
					"api_method":           operationName,
					"api_name":             apiName,
					"api_version":          "1.0.0",
					"http_status_code":     status,
					"location":             "us-west1",
					"log_message":          operationName + " is called",
					"producer_project_id":  "nerdvana",
					"response_code_detail": "via_upstream",
					"service_agent":        "ESPv2/2.45.0",
					"service_config_id":    serviceConfig,
					"timestamp":            timestampfloat,
					"xtra":                 "extra info",
				}
				pbytes, err := json.Marshal(payload)
				if err != nil {
					return err
				}
				operation.EndTime = timestamp2
				operation.LogEntries = []*servicecontrol.LogEntry{
					{
						Name:      "endpoints_log",
						Timestamp: timestamp,
						Severity:  "INFO",
						HttpRequest: &servicecontrol.HttpRequest{
							RequestMethod: "GET",
							RequestUrl:    "/unknown",
							RequestSize:   10,
							Status:        int64(status),
							ResponseSize:  10,
							RemoteIp:      callerIP,
							Latency:       "10s",
							Protocol:      "http",
						},
						StructPayload: pbytes,
					},
				}
				operation.MetricValueSets = []*servicecontrol.MetricValueSet{
					createInt64MetricSet("serviceruntime.googleapis.com/api/consumer/request_count", 1),
					createInt64MetricSet("serviceruntime.googleapis.com/api/producer/request_count", 1),
					createInt64MetricSet("serviceruntime.googleapis.com/api/consumer/quota_used_count", 1),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/consumer/total_latencies", 1),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/total_latencies", 1),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/request_sizes", 200),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/response_sizes", 200),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/request_overhead_latencies", 1),
					createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/backend_latencies", 1),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/request_sizes", 200),
					createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/response_sizes", 200),
				}
				operation.QuotaProperties = nil
				request := &servicecontrol.ReportRequest{
					Operations: []*servicecontrol.Operation{
						operation,
					},
				}

				start := time.Now()
				result, err := srv.Services.Report(serviceName, request).Do()
				if err != nil {
					return err
				}
				elapsed := time.Since(start)

				bytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return &json.UnsupportedValueError{}
				}
				fmt.Printf("REPORT (%fms) %d > %s\n", float64(elapsed)/1e6, i, string(bytes))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of: (yaml, json).")
	cmd.Flags().IntVarP(&iterations, "iterations", "i", 1, "Number of times to call report.")
	return cmd
}
