#!/bin/bash
set -e

REGION=$(gcloud config get run/region)
SPEC_REVISION=$(registry get projects/timburks-test/locations/global/apis/bookstore/deployments/backend | jq .apiSpecRevision -r)
ENDPOINT_ADDRESS=https://$(gcloud api-gateway gateways describe bookstore --location us-west2 --format "value(defaultHostname)")

cat > gateway-deployment.yaml <<EOF
apiVersion: apigeeregistry/v1
kind: Deployment
metadata:
  name: gateway
  parent: apis/bookstore
  labels:
    platform: apigateway
    apihub-gateway: apihub-google-cloud-api-gateway
  annotations:
    region: $REGION
data:
  displayName: Gateway
  description: An API Gateway deployment of the Bookstore API
  apiSpecRevision: $SPEC_REVISION
  endpointURI: $ENDPOINT_ADDRESS
  intendedAudience: Public
EOF
