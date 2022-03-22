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


# Build this image using "envsubst < containers/registry-tools/registry-linters/Dockerfile | docker build -t gcr.io/$REGISTRY_PROJECT_IDENTIFIER/registry-linters -f - ."
# docker push gcr.io/$REGISTRY_PROJECT_IDENTIFIER/registry-linters:latest
# Make sure REGISTRY_PROJECT_IDENTIFIER is set

# Build stage
FROM ghcr.io/apigee/registry-tools:latest as registry-tools

FROM gcr.io/cloud-builders/kubectl
ENV MOCKSERVICE_DOMAIN=""
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
ENV APG_REGISTRY_INSECURE=0

COPY --from=registry-tools /bin/registry /bin/registry
COPY --from=registry-tools /bin/apg /bin/apg

COPY ./create-mock.sh /create-mock.sh
COPY ./docker-entrypoint.sh /docker-entrypoint.sh

COPY ./mock-server-deployment.template.yaml /mock-server-deployment.template.yaml

RUN apt-get update -y && apt-get install gettext-base -y \
    && chmod +x /create-mock.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]