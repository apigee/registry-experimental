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

FROM ghcr.io/apigee/registry-tools:latest as registry-tools
ENV REGISTRY_SPEC=""
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
ENV APG_REGISTRY_INSECURE=0

FROM launcher.gcr.io/google/nodejs
ENV REGISTRY_SPEC=""
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
ENV APG_REGISTRY_INSECURE=0

EXPOSE 4010

COPY --from=registry-tools /bin/registry /bin/registry
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && npm install -g @stoplight/prism-cli

ENTRYPOINT ["/entrypoint.sh"]

