# Copyright 2023 Google LLC
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

# This Dockerfile builds an image that contains the registry and registry-experimental tools
# along with several associated API tools including linters and code generators.

# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.20.1 as builder
RUN apt-get update

### BUILD REGISTRY AND PLUGINS ###

# Get the latest registry source code
WORKDIR /
RUN git clone https://github.com/apigee/registry

# Work in the registry directory
WORKDIR /registry

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
RUN go mod download

# Build the registry tool and plugins.
RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/registry
RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/registry/plugins/registry-lint-api-linter
RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/registry/plugins/registry-lint-spectral

### BUILD REGISTRY-EXPERIMENTAL ###

# Get the latest registry-experimental source code
WORKDIR /
RUN git clone https://github.com/apigee/registry-experimental

# Work in the registry-experimental directory
WORKDIR /registry-experimental

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
RUN go mod download

# Build registry-experimental.
RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/registry-experimental

### Build other Go-based tools ###
WORKDIR /
RUN CGO_ENABLED=0 GOOS=linux go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
RUN CGO_ENABLED=0 GOOS=linux go install github.com/googleapis/api-linter/cmd/api-linter@latest

### BUILD THE CONTAINER ###

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk update && apk upgrade
RUN apk add --no-cache ca-certificates

# Add some commonly-needed tools to the image.
RUN apk add --no-cache bash git yq jq

# Copy registry binaries from the builder stage.
COPY --from=builder /registry/registry /usr/local/bin/registry
COPY --from=builder /registry/registry-lint-api-linter /usr/local/bin/registry-lint-api-linter
COPY --from=builder /registry/registry-lint-spectral /usr/local/bin/registry-lint-spectral
COPY --from=builder /registry-experimental/registry-experimental /usr/local/bin/registry-experimental

# Copy installed Go tools from the builder stage.
COPY --from=builder /go/bin/protoc-gen-openapi /usr/local/bin/protoc-gen-openapi
COPY --from=builder /go/bin/api-linter /usr/local/bin/api-linter

# Install protoc.
RUN apk add curl gcompat
COPY --from=builder /registry/tools/PROTOC-VERSION.sh ./tools/PROTOC-VERSION.sh
COPY --from=builder /registry/tools/FETCH-PROTOC.sh ./tools/FETCH-PROTOC.sh
RUN ./tools/FETCH-PROTOC.sh

# Install node-based tools.
RUN apk add --no-cache nodejs npm
RUN npm install -g @stoplight/spectral-cli
RUN npm install -g api-spec-converter