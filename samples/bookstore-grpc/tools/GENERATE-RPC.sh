#!/bin/bash
set -e

. tools/PROTOS.sh

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

for proto in ${SERVICE_PROTOS[@]}; do
	echo "Generating Go types for $proto"
	protoc $proto --proto_path='.' \
		--proto_path=$COMMON_PROTOS_PATH \
		--go_opt='module=github.com/examples/bookstore/rpc' \
		--go_out='./rpc'
done
