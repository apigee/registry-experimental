#!/bin/bash
set -e

. tools/PROTOS.sh

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Generating Go gRPC client/server for ${SERVICE_PROTOS[@]}"
protoc ${SERVICE_PROTOS[*]} --proto_path='.' \
  --proto_path=$COMMON_PROTOS_PATH \
  --go-grpc_opt='module=github.com/examples/bookstore' \
  --go-grpc_out='.'
