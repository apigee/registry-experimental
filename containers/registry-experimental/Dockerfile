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

# This Dockerfile builds an image that runs the authz-server.

# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18-rc as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build registry-experimental.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o registry-experimental ./cmd/registry-experimental

# Get registry code.
RUN git clone https://github.com/apigee/registry

# Build registry tool and apg.
RUN cd registry; CGO_ENABLED=0 GOOS=linux go build -v -o registry ./cmd/registry
RUN cd registry; CGO_ENABLED=0 GOOS=linux go build -v -o apg ./cmd/apg

# Install protoc-gen-openapi plugin
RUN go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

###
### Note that all binaries are installed to /usr/local/bin
###

# Copy the registry-experimental binary from the builder stage.
COPY --from=builder /app/registry-experimental /usr/local/bin/registry-experimental

# Copy the registry and apg binaries from the builder stage.
COPY --from=builder /app/registry/registry /usr/local/bin/registry
COPY --from=builder /app/registry/apg /usr/local/bin/apg

# Install node-based tools
RUN apk add --no-cache nodejs npm
RUN npm install -g @stoplight/spectral@5.9.2
RUN npm install -g api-spec-converter

# Install protoc
RUN apk add curl gcompat
COPY --from=builder /app/registry/tools/FETCH-PROTOC.sh ./FETCH-PROTOC.sh
RUN ./FETCH-PROTOC.sh

# Copy protoc-gen-openapi from the builder stage.
COPY --from=builder /go/bin/protoc-gen-openapi /usr/local/bin/protoc-gen-openapi
