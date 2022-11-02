#!/bin/bash
set -e

. tools/PROTOS.sh

go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

echo "Generating OpenAPI spec for ${SERVICE_PROTOS[@]}"
protoc ${SERVICE_PROTOS[*]} --proto_path='.' --proto_path=$COMMON_PROTOS_PATH --openapi_out='.'
