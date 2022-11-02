#!/bin/bash
set -e

. tools/PROTOS.sh

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_cli@latest
go install golang.org/x/tools/cmd/goimports@latest

# This directory contains the generated CLI
GENERATED='cmd/bookstore'

echo "Generating Go client CLI for ${SERVICE_PROTOS[@]}"
protoc ${SERVICE_PROTOS[*]} \
	--proto_path='.' \
	--proto_path=$COMMON_PROTOS_PATH \
  	--go_cli_opt='root=bookstore' \
  	--go_cli_opt='gapic=github.com/examples/bookstore/gapic' \
	--go_cli_out=$GENERATED

goimports -w cmd/bookstore