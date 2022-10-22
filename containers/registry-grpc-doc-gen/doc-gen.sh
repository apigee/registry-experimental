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

mkdir -p /workspace/args[0]/protos


TOKEN=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | jq .access_token -r)
SPEC_DETAILS=$(registry get args[0] --token=$TOKEN --registry.address=${REGISTRY_ADDRESS})
MIMETYPE=$( echo $SPEC_DETAILS | jq -r .mimeType)
FILENAME=$( echo $SPEC_DETAILS | jq -r .filename)

SPEC_PATH="/workspace/${args[0]}"

registry get arg[0] --contents --token=$TOKEN --registry.address=${REGISTRY_ADDRESS} > $SPEC_PATH/$FILENAME

if [$MIMETYPE eq "application/x.protobuf+gzip"]
then
  tar xvfz "$SPEC_PATH/$FILENAME" "$SPEC_PATH/protos"
fi

protoc "$SPEC_PATH/**/*.proto" --proto_path="/googleapis-common-protos" --doc_out="$SPEC_PATH" --doc_opt=html,index.html

registry rpc create-artifact grpc-documentation \
  --artifact.contents=$(cat "$SPEC_PATH/index.html" | od -A n -t x1 | sed 's/ *//g') \
  --artifact.mime_type="text/html" \
  --artifact_id=grpc-documentation \
  --parent=args[0]

