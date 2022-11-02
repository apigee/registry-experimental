#!/bin/bash
set -e

PROJECT=$(gcloud config get project)
ADDRESS=$(registry get apis/bookstore/deployments/backend | jq .endpointUri -r)
ADDRESS=${ADDRESS#https://}

cat > api_config.yaml <<EOF
# The configuration schema is defined by the service.proto file.
# https://github.com/googleapis/googleapis/blob/master/google/api/service.proto

type: google.api.Service
config_version: 3
name: "*.apigateway.$PROJECT.cloud.goog"
title: API Gateway + Cloud Run gRPC
apis:
  - name: examples.bookstore.v1.Bookstore
usage:
  rules:
  - selector: "*"
    allow_unregistered_calls: true
backend:
  rules:
  - selector: "*"
    address: grpcs://$ADDRESS
EOF