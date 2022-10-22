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

FROM golang:1.18 as protoc-gen-doc
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

FROM ghcr.io/apigee/registry-tools:latest

COPY --from=protoc-gen-doc /go/bin/protoc-gen-doc /bin/

ARG PB_VERSION="21.8"
RUN PB_REL="https://github.com/protocolbuffers/protobuf/releases" && \
    curl -LO ${PB_REL}/download/v${PB_VERSION}/protoc-${PB_VERSION}-linux-x86_64.zip && \
    mkdir /protoc && \
    unzip protoc-${PB_VERSION}-linux-x86_64.zip -d /protoc && \
    ln -s /protoc/bin/protoc /bin/protoc


COPY ./doc-gen.sh /
RUN chmod +x /doc-gen.sh

RUN git clone https://github.com/googleapis/api-common-protos /googleapis-common-protos