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

package servicebuilder

import (
	"fmt"
	"strings"

	"google.golang.org/api/servicemanagement/v1"
)

type Operation struct {
	Id     string
	Method string
	Path   string
}

type Api struct {
	Name       string
	Version    string
	Operations []Operation
}

func Apis(apis []*Api) []*servicemanagement.Api {
	response := []*servicemanagement.Api{}
	for i, a := range apis {
		apiName := fmt.Sprintf("%d.%s", i+1, strings.ReplaceAll(a.Name, ".", "_"))
		methods := []*servicemanagement.Method{}
		for _, op := range a.Operations {
			methods = append(methods, &servicemanagement.Method{
				Name: op.Id,
				// config creation fails (with 500) if the following two fields are omitted
				RequestTypeUrl:  "type.googleapis.com/google.protobuf.Empty",
				ResponseTypeUrl: "type.googleapis.com/google.protobuf.Value",
			})
		}
		response = append(response, &servicemanagement.Api{
			Name:    apiName,
			Methods: methods,
			Version: a.Version,
		})
	}
	return response
}

func Control() *servicemanagement.Control {
	return &servicemanagement.Control{
		Environment: "servicecontrol.googleapis.com",
	}
}

func Enums() []*servicemanagement.Enum {
	return []*servicemanagement.Enum{
		{
			Enumvalue: []*servicemanagement.EnumValue{
				{
					Name: "NULL_VALUE",
				},
			},
			Name: "google.protobuf.NullValue",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
		},
	}
}

func Http(apis []*Api) *servicemanagement.Http {
	rules := []*servicemanagement.HttpRule{}
	for i, a := range apis {
		apiName := fmt.Sprintf("%d.%s", i+1, strings.ReplaceAll(a.Name, ".", "_"))
		for _, op := range a.Operations {
			switch op.Method {
			case "GET":
				rules = append(rules, &servicemanagement.HttpRule{
					Get:      op.Path,
					Selector: apiName + "." + op.Id,
				})
			case "POST":
				rules = append(rules, &servicemanagement.HttpRule{
					Post:     op.Path,
					Selector: apiName + "." + op.Id,
				})
			default:
				panic(op.Method)
			}
		}
	}
	return &servicemanagement.Http{
		Rules: rules,
	}
}
func Logging() *servicemanagement.Logging {
	return &servicemanagement.Logging{
		ProducerDestinations: []*servicemanagement.LoggingDestination{
			{
				Logs:              []string{"endpoints_log"},
				MonitoredResource: "api",
			},
		},
	}
}

func Logs() []*servicemanagement.LogDescriptor {
	return []*servicemanagement.LogDescriptor{
		{
			Name: "endpoints_log",
		},
	}
}

func Metrics() []*servicemanagement.MetricDescriptor {
	return []*servicemanagement.MetricDescriptor{
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/credential_id"},
				{Key: "/protocol"},
				{Key: "/response_code"},
				{Key: "/response_code_class"},
				{Key: "/status_code"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/consumer/request_count",
			Type:       "serviceruntime.googleapis.com/api/consumer/request_count",
			ValueType:  "INT64",
		},
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/consumer/total_latencies",
			Type:       "serviceruntime.googleapis.com/api/consumer/total_latencies",
			ValueType:  "DISTRIBUTION",
		},
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/protocol"},
				{Key: "/response_code"},
				{Key: "/response_code_class"},
				{Key: "/status_code"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/request_count",
			Type:       "serviceruntime.googleapis.com/api/producer/request_count",
			ValueType:  "INT64",
		},
		{
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/total_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/total_latencies",
			ValueType:  "DISTRIBUTION",
		},
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/credential_id"},
				{Key: "/quota_group_name"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/consumer/quota_used_count",
			Type:       "serviceruntime.googleapis.com/api/consumer/quota_used_count",
			ValueType:  "INT64",
		},
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/consumer/request_sizes",
			Type:       "serviceruntime.googleapis.com/api/consumer/request_sizes",
			ValueType:  "DISTRIBUTION",
		},
		{
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "/credential_id"},
			},
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/consumer/response_sizes",
			Type:       "serviceruntime.googleapis.com/api/consumer/response_sizes",
			ValueType:  "DISTRIBUTION",
		},
		{
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
			ValueType:  "DISTRIBUTION",
		},
		{
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/backend_latencies",
			Type:       "serviceruntime.googleapis.com/api/producer/backend_latencies",
			ValueType:  "DISTRIBUTION",
		},
		{
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/request_sizes",
			Type:       "serviceruntime.googleapis.com/api/producer/request_sizes",
			ValueType:  "DISTRIBUTION",
		},
		{
			MetricKind: "DELTA",
			Name:       "serviceruntime.googleapis.com/api/producer/response_sizes",
			Type:       "serviceruntime.googleapis.com/api/producer/response_sizes",
			ValueType:  "DISTRIBUTION",
		},
		{
			Name:        "calls",
			DisplayName: "Calls",
			ValueType:   "INT64",
			MetricKind:  "DELTA",
		},
	}
}

func MonitoredResources() []*servicemanagement.MonitoredResourceDescriptor {
	return []*servicemanagement.MonitoredResourceDescriptor{
		{
			Type: "api",
			Labels: []*servicemanagement.LabelDescriptor{
				{Key: "cloud.googleapis.com/location"},
				{Key: "cloud.googleapis.com/uid"},
				{Key: "serviceruntime.googleapis.com/api_version"},
				{Key: "serviceruntime.googleapis.com/api_method"},
				{Key: "serviceruntime.googleapis.com/consumer_project"},
				{Key: "cloud.googleapis.com/project"},
				{Key: "cloud.googleapis.com/service"},
			},
		},
	}
}

func Monitoring() *servicemanagement.Monitoring {
	return &servicemanagement.Monitoring{
		ConsumerDestinations: []*servicemanagement.MonitoringDestination{
			{
				Metrics: []string{
					"serviceruntime.googleapis.com/api/consumer/request_count",
					"serviceruntime.googleapis.com/api/consumer/quota_used_count",
					"serviceruntime.googleapis.com/api/consumer/total_latencies",
					"serviceruntime.googleapis.com/api/consumer/request_sizes",
					"serviceruntime.googleapis.com/api/consumer/response_sizes",
				},
				MonitoredResource: "api",
			},
		},
		ProducerDestinations: []*servicemanagement.MonitoringDestination{
			{
				Metrics: []string{
					"serviceruntime.googleapis.com/api/producer/request_count",
					"serviceruntime.googleapis.com/api/producer/total_latencies",
					"serviceruntime.googleapis.com/api/producer/request_overhead_latencies",
					"serviceruntime.googleapis.com/api/producer/backend_latencies",
					"serviceruntime.googleapis.com/api/producer/request_sizes",
					"serviceruntime.googleapis.com/api/producer/response_sizes",
				},
				MonitoredResource: "api",
			},
		},
	}
}

func Types() []*servicemanagement.Type {
	return []*servicemanagement.Type{
		{
			Name: "google.protobuf.ListValue",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
			Fields: []*servicemanagement.Field{
				{
					Cardinality: "CARDINALITY_REPEATED",
					JsonName:    "values",
					Kind:        "TYPE_MESSAGE",
					Name:        "values",
					Number:      1,
					TypeUrl:     "type.googleapis.com/google.protobuf.Value",
				},
			},
		},
		{
			Name: "google.protobuf.Struct",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
			Fields: []*servicemanagement.Field{
				{
					Cardinality: "CARDINALITY_REPEATED",
					JsonName:    "fields",
					Kind:        "TYPE_MESSAGE",
					Name:        "fields",
					Number:      1,
					TypeUrl:     "type.googleapis.com/google.protobuf.Struct.FieldsEntry",
				},
			},
		},
		{
			Name: "google.protobuf.Struct.FieldsEntry",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
			Fields: []*servicemanagement.Field{
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "key",
					Kind:        "TYPE_STRING",
					Name:        "key",
					Number:      1,
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "value",
					Kind:        "TYPE_MESSAGE",
					Name:        "value",
					Number:      2,
					TypeUrl:     "type.googleapis.com/google.protobuf.Value",
				},
			},
		},
		{
			Name: "google.protobuf.Empty",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
		},
		{
			Name: "google.protobuf.Value",
			SourceContext: &servicemanagement.SourceContext{
				FileName: "struct.proto",
			},
			Fields: []*servicemanagement.Field{
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "nullValue",
					Kind:        "TYPE_ENUM",
					Name:        "null_value",
					Number:      1,
					TypeUrl:     "type.googleapis.com/google.protobuf.NullValue",
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "numberValue",
					Kind:        "TYPE_DOUBLE",
					Name:        "number_value",
					Number:      2,
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "stringValue",
					Kind:        "TYPE_STRING",
					Name:        "string_value",
					Number:      3,
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "boolValue",
					Kind:        "TYPE_BOOL",
					Name:        "bool_value",
					Number:      4,
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "structValue",
					Kind:        "TYPE_MESSAGE",
					Name:        "struct_value",
					Number:      5,
					TypeUrl:     "type.googleapis.com/google.protobuf.Struct",
				},
				{
					Cardinality: "CARDINALITY_OPTIONAL",
					JsonName:    "listValue",
					Kind:        "TYPE_MESSAGE",
					Name:        "list_value",
					Number:      6,
					TypeUrl:     "type.googleapis.com/google.protobuf.ListValue",
				},
			},
		},
	}
}

func Usage(apis []*Api) *servicemanagement.Usage {
	rules := []*servicemanagement.UsageRule{}
	for i, a := range apis {
		apiName := fmt.Sprintf("%d.%s", i+1, strings.ReplaceAll(a.Name, ".", "_"))
		for _, op := range a.Operations {
			rules = append(rules, &servicemanagement.UsageRule{
				Selector: apiName + "." + op.Id,
			})
		}
	}
	return &servicemanagement.Usage{
		Rules: rules,
	}
}
