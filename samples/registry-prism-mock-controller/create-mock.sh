#!/bin/bash
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -ex

export REGISTRY_SPEC=$1

# REGISTRY_SPEC is in  format projects/p1/locations/global/apis/ap1/versions/v1/specs/spec1
echo $REGISTRY_SPEC | cut -d "/" -f6
API="$(echo $REGISTRY_SPEC | cut -d "/" -f6)"
VERSION="$(echo $REGISTRY_SPEC | cut -d "/" -f8)"
SPEC="$(echo $REGISTRY_SPEC | cut -d "/" -f10)"

export REGISTRY_VERSION_SPEC="$(echo $API-$VERSION-$SPEC | sed 's/[^a-zA-Z0-9\-]/-/g')"
export MOCK_SERVICE_ENDPOINT="mock-${REGISTRY_VERSION_SPEC}.${MOCKSERVICE_DOMAIN}"

cat /mock-server-deployment.template.yaml | envsubst > /tmp/mock-server-deployment-${REGISTRY_VERSION_SPEC}.yaml

kubectl apply -f /tmp/mock-server-deployment-${REGISTRY_VERSION_SPEC}.yaml

apg registry create-artifact \
  --parent $REGISTRY_SPEC \
  --artifact_id "prism-mock-endpoint" \
  --artifact.mime_type "text/plain" \
  --artifact.contents $(echo $MOCK_SERVICE_ENDPOINT | od -A n -t x1) \
  --json
