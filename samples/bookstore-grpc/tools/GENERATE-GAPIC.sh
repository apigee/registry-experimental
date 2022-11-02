#!/bin/bash
set -e

. tools/PROTOS.sh

go install github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic@latest
go install golang.org/x/tools/cmd/goimports@latest

echo "Generating Go client library for ${SERVICE_PROTOS[@]}"
protoc ${SERVICE_PROTOS[*]} \
	--proto_path='.' \
	--proto_path=$COMMON_PROTOS_PATH \
	--go_gapic_opt='go-gapic-package=github.com/examples/bookstore/gapic;gapic' \
	--go_gapic_opt='grpc-service-config=gapic/grpc_service_config.json' \
	--go_gapic_opt='module=github.com/examples/bookstore' \
	--go_gapic_out='.'

goimports -w gapic