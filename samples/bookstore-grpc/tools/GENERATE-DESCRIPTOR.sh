#!/bin/bash
set -e

#. tools/PROTOS.sh
#
#echo "Generating descriptor for ${SERVICE_PROTOS[@]}"
#protoc ${SERVICE_PROTOS[*]} \
#	--proto_path='.' \
#	--proto_path=$COMMON_PROTOS_PATH \
#	--descriptor_set_out='proto.pb' \
#	--include_imports \
#	--include_source_info

registry get projects/timburks-test/locations/global/apis/bookstore/versions/v1/specs/examples-bookstore-v1.zip/artifacts/descriptor --raw > proto.pb
