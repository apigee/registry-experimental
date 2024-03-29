# This Dockerfile builds an image that runs the registry-server with no proxy.

# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18 as builder
# Create and change to the app directory.
WORKDIR /app
# Copy local code to the container image.
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o dispatcher ./cmd/workers/dispatcher-server/main.go

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk update
RUN apk add bash
RUN apk add --no-cache ca-certificates
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/dispatcher /dispatcher
CMD ["/dispatcher"]
