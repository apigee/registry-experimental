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
	"net/http"
	"strings"
	"time"

	"github.com/apigee/registry-experimental/cmd/zero/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/servicecontrol/v1"
)

func Middleware(serviceName string, verbose bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		t, err := NewTracker(c, serviceName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		t.Verbose = verbose
		err = t.Check()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		err = t.AllocateQuota()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		t.CallHandler()
		err = t.Report()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

type Tracker struct {
	Context         *gin.Context
	ServiceName     string
	Client          *http.Client
	Service         *servicecontrol.Service
	StartTime       time.Time
	BackendDuration time.Duration
	Operation       *servicecontrol.Operation
	Config          *config.Config
	Method          string
	ApiKey          string
	Verbose         bool
}

func NewTracker(gc *gin.Context, serviceName string) (*Tracker, error) {
	apiName := "1." + strings.ReplaceAll(
		strings.ReplaceAll(serviceName, ".", "_"),
		"-", "_")
	// this assumes/requires that the handler function name exactly matches the operation name
	parts := strings.Split(gc.HandlerName(), ".")
	method := parts[len(parts)-1]
	var err error
	tracker := &Tracker{
		Context:     gc,
		ServiceName: serviceName,
		Method:      method,
	}
	ctx := context.Background()
	tracker.Client, err = config.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	tracker.Service, err = servicecontrol.NewService(ctx, option.WithHTTPClient(tracker.Client))
	if err != nil {
		return nil, err
	}
	tracker.ApiKey = gc.Request.Header.Get("X-Api-Key")
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	tracker.Config = c
	tracker.StartTime = time.Now()
	tracker.Operation = &servicecontrol.Operation{
		OperationId:   uuid.New().String(),
		OperationName: apiName + "." + method,
		ConsumerId:    "api_key:" + tracker.ApiKey,
		StartTime:     tracker.StartTime.Format(time.RFC3339Nano),
	}
	return tracker, nil
}

func (t *Tracker) Check() error {
	request := &servicecontrol.CheckRequest{
		Operation: t.Operation,
	}
	start := time.Now()
	result, err := t.Service.Services.Check(t.ServiceName, request).Do()
	elapsed := time.Since(start)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &json.UnsupportedValueError{}
	}
	if t.Verbose {
		fmt.Printf("CHECK (%fms)> %s\n", float64(elapsed)/1e6, string(bytes))
	}
	if result.CheckErrors != nil {
		return fmt.Errorf("%s", result.CheckErrors[0].Code)
	}
	return nil
}

func (t *Tracker) AllocateQuota() error {
	apiName := "1." + strings.ReplaceAll(
		strings.ReplaceAll(t.ServiceName, ".", "_"),
		"-", "_")
	operationName := apiName + "." + t.Method
	request := &servicecontrol.AllocateQuotaRequest{
		ServiceConfigId: t.Config.ServiceConfig,
		AllocateOperation: &servicecontrol.QuotaOperation{
			ConsumerId:  "api_key:" + t.ApiKey,
			OperationId: t.Operation.OperationId,
			QuotaMode:   "NORMAL",
			MethodName:  operationName,
			QuotaMetrics: []*servicecontrol.MetricValueSet{
				createInt64MetricSet("calls", 1),
			},
		},
	}
	start := time.Now()
	result, err := t.Service.Services.AllocateQuota(t.ServiceName, request).Do()
	elapsed := time.Since(start)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &json.UnsupportedValueError{}
	}
	if t.Verbose {
		fmt.Printf("ALLOCATE QUOTA (%fms) > %s\n", float64(elapsed)/1e6, string(bytes))
	}
	if result.AllocateErrors != nil {
		return errors.New("out of quota")
	}
	return nil
}

func (t *Tracker) CallHandler() {
	start := time.Now()
	t.Context.Next()
	t.BackendDuration = time.Since(start)
}

func (t *Tracker) Report() error {
	status := t.Context.Writer.Status()
	now := time.Now()
	timestampfloat := float64(now.UnixNano()) / 1e9
	timestamp2 := now.Format(time.RFC3339Nano)

	latency := time.Since(t.StartTime)

	requestSize := t.Context.Request.ContentLength
	responseSize := t.Context.Writer.Size()

	callerIP := t.Context.RemoteIP()
	apiName := "1." + strings.ReplaceAll(
		strings.ReplaceAll(t.ServiceName, ".", "_"),
		"-", "_")
	operationName := apiName + "." + t.Method
	operation := t.Operation
	payload := map[string]interface{}{
		"api_key_state":        "CHECKED",
		"api_key":              t.ApiKey,
		"api_method":           operationName,
		"api_name":             apiName,
		"api_version":          "1.0.0",
		"http_status_code":     status,
		"location":             "local",
		"log_message":          operationName + " is called",
		"producer_project_id":  t.Config.ProducerProject,
		"response_code_detail": "via_upstream",
		"service_agent":        "Zero/0.0.1",
		"service_config_id":    t.Config.ServiceConfig,
		"timestamp":            timestampfloat,
		"xtra":                 "extra info",
	}
	pbytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	producerProject := t.Config.ProducerProject
	operation.EndTime = timestamp2
	operation.Labels = map[string]string{
		"cloud.googleapis.com/location":             "us-west1",
		"serviceruntime.googleapis.com/api_method":  operationName,
		"cloud.googleapis.com/project":              producerProject,
		"cloud.googleapis.com/service":              t.ServiceName,
		"serviceruntime.googleapis.com/api_version": "1.0.0",
		// none of the following appear in the logs
		// but they seem to be set in ESPv2
		//  https://github.com/GoogleCloudPlatform/esp-v2/blob/9217d68484321aceb7f7fbdc63be9363c96ed722/tests/utils/service_control_utils.go#L245
		"cloud.googleapis.com/uid":                       t.Operation.OperationId,
		"servicecontrol.googleapis.com/caller_ip":        callerIP,
		"servicecontrol.googleapis.com/service_agent":    "Zero/0.0.1",
		"servicecontrol.googleapis.com/platform":         "Custom",
		"servicecontrol.googleapis.com/user_agent":       "Zero",
		"serviceruntime.googleapis.com/consumer_project": t.Config.ConsumerProject,
		"/response_code":                                 fmt.Sprintf("%d", status),
		"/response_code_class":                           fmt.Sprintf("%dxx", status/100),
		"/status_code":                                   fmt.Sprintf("%d", status),
		"/protocol":                                      "http",
	}
	operation.LogEntries = []*servicecontrol.LogEntry{
		{
			Name:      "endpoints_log",
			Timestamp: timestamp2,
			Severity:  "INFO",
			HttpRequest: &servicecontrol.HttpRequest{
				RequestMethod: t.Context.Request.Method,
				RequestUrl:    t.Context.Request.URL.Path,
				RequestSize:   requestSize,
				Status:        int64(status),
				ResponseSize:  int64(responseSize),
				RemoteIp:      callerIP,
				Latency:       fmt.Sprintf("%fs", latency.Seconds()),
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
		createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/request_sizes", requestSize),
		createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/consumer/response_sizes", int64(responseSize)),
		createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/request_overhead_latencies", 1),
		createDistMetricSet(&timeDistOptions, "serviceruntime.googleapis.com/api/producer/backend_latencies", t.BackendDuration.Milliseconds()),
		createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/request_sizes", requestSize),
		createDistMetricSet(&sizeDistOptions, "serviceruntime.googleapis.com/api/producer/response_sizes", int64(responseSize)),
	}
	operation.QuotaProperties = nil
	request := &servicecontrol.ReportRequest{
		Operations: []*servicecontrol.Operation{
			operation,
		},
	}
	start := time.Now()
	result, err := t.Service.Services.Report(t.ServiceName, request).Do()
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &json.UnsupportedValueError{}
	}
	if t.Verbose {
		fmt.Printf("REPORT (%fms) > %s\n", float64(elapsed)/1e6, string(bytes))
	}
	return nil
}

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
