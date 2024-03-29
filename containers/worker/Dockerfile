# This Dockerfile builds an image that runs the registry-server with no proxy.

# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18 as builder
RUN apt-get update
RUN apt-get install unzip
RUN curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip > /tmp/protoc.zip
RUN unzip /tmp/protoc.zip -d /usr/local
# Create and change to the app directory.
WORKDIR /app
# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download
# Copy local code to the container image.
COPY . ./
# Build the worker
RUN CGO_ENABLED=0 GOOS=linux go build -v -o worker ./cmd/workers/worker-server/main.go

# Build the registry-tool.
RUN git clone https://github.com/apigee/registry
RUN cd registry; CGO_ENABLED=0 GOOS=linux go build -v ./cmd/registry

# Install the api-linter
RUN go get -u github.com/googleapis/api-linter/cmd/api-linter

# Install the gnostic linter
RUN git clone https://github.com/google/gnostic
RUN cd gnostic; go install ./...

# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
# We are using debian base image here because the api-linter tool build on the builder doesn't work with alpine image.
FROM debian:latest
RUN apt-get update
RUN apt-get install -y git ca-certificates
# Copy the linter binaries
COPY --from=builder /go/bin/api-linter /bin/
COPY --from=builder /go/bin/gnostic* /bin/
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/worker /worker
# Copy the registry tool binary
COPY --from=builder /app/registry /bin/
ENV PATH=$PATH:/bin/
CMD ["/worker"]
