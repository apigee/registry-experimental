#!/bin/bash

# Copyright 2022 Google LLC
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
set -euox pipefail


args=("$@")
SPEC=${args[0]}
TOKEN=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | jq .access_token -r)
SPEC_PATH="/tmp/workspace/$SPEC"

rm -rf "$SPEC_PATH"

mkdir -p "$SPEC_PATH/protos"

SPEC_DETAILS=$(registry get  $SPEC --registry.token=$TOKEN --registry.address=${REGISTRY_ADDRESS})
MIMETYPE=$( echo $SPEC_DETAILS | jq -r .mimeType)
FILENAME=$( echo $SPEC_DETAILS | jq -r .filename)

registry get $SPEC  --contents --registry.token=$TOKEN --registry.address=${REGISTRY_ADDRESS} > "$SPEC_PATH/$FILENAME"

if [ $MIMETYPE == "application/x.protobuf+gzip" ]; then
  tar -xf "$SPEC_PATH/$FILENAME" -C "$SPEC_PATH/protos"
fi

find "$SPEC_PATH/protos" -type f -name "*.proto"  > "$SPEC_PATH/proto-files.txt"

protoc @"$SPEC_PATH/proto-files.txt" \
  --proto_path="$SPEC_PATH/protos" \
  --proto_path="/protoc/include" \
  --proto_path="/tmp/googleapis-common-protos" \
  --doc_out="$SPEC_PATH" --doc_opt=html,index.html

DOC_HEX_STRING=$(cat "$SPEC_PATH/index.html" | xxd -ps -c 200 | tr -d '\n')

registry rpc create-artifact grpc-documentation \
  --artifact.mime_type "text/plain" \
  --artifact_id grpc-documentation \
  --parent $SPEC \
  --artifact.contents "$DOC_HEX_STRING"


