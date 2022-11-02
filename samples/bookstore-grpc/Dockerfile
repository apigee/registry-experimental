FROM golang:1.18.1 as builder
RUN apt-get update
RUN apt-get install unzip

WORKDIR /work

COPY go.mod ./
RUN go mod download

COPY . .

RUN mkdir -p ./rpc
RUN ./tools/FETCH-PROTOC.sh
RUN ./tools/GENERATE-RPC.sh
RUN ./tools/GENERATE-GRPC.sh

# Build bookstore-server.
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN goimports -w ./cmd/bookstore-server
RUN go get ./cmd/bookstore-server
RUN CGO_ENABLED=0 GOOS=linux go build -v -o bookstore-server ./cmd/bookstore-server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /work/bookstore-server /bookstore-server

# Run the service on container startup.
CMD ["/bookstore-server"]
