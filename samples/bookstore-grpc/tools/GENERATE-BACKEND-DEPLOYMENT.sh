#!/bin/bash
set -e

REGION=$(gcloud config get run/region)
REVISION=$(registry get projects/timburks-test/locations/global/apis/bookstore/versions/v1/specs/examples-bookstore-v1.zip | jq .revisionId -r)
ADDRESS=$(gcloud run services describe bookstore --format "value(status.url)")

cat > backend-deployment.yaml <<EOF
apiVersion: apigeeregistry/v1
kind: Deployment
metadata:
  name: backend
  parent: apis/bookstore
  labels:
    platform: cloudrun
    apihub-gateway: apihub-unmanaged
  annotations:
    region: $REGION
data:
  displayName: Backend
  description: The backend deployment of the Bookstore API
  apiSpecRevision: v1/specs/examples-bookstore-v1.zip@$REVISION
  endpointURI: $ADDRESS
  intendedAudience: Internal
EOF
