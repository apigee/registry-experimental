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
set -euo pipefail

# command expect name of the spec and the name of the bucket where html docs will be generated
args=("$@")
SPEC=${args[0]}
BUCKET_NAME=${args[1]}
PROTO_PATH=${args[2]}

if [[ ${TOKEN:-"unset"} == "unset" ]]; then
  TOKEN=$(cat /tmp/registry-token)
fi

GRPC_GEN_DOC_FILE="index.html"

echo "Started processing $SPEC"

SPEC_PATH="/tmp/workspace/$SPEC"

rm -rf "$SPEC_PATH"

mkdir -p "$SPEC_PATH"

SPEC_DETAILS=$(registry get  "$SPEC")
MIMETYPE=$( echo "$SPEC_DETAILS" | jq -r .mimeType)
FILENAME=$( echo "$SPEC_DETAILS" | jq -r .filename)

echo "About to get contents of $SPEC"

registry get "$SPEC"  \
    --contents > "$SPEC_PATH/$FILENAME"

if [ "$MIMETYPE" == "application/x.protobuf+gzip" ]; then
  tar -xf "$SPEC_PATH/$FILENAME" -C "$SPEC_PATH"
  # Mac OS gziped files may contain archive files
  # https://en.wikipedia.org/wiki/AppleSingle_and_AppleDouble_formats
  rm -rf "$SPEC_PATH/.*"
fi

find "$SPEC_PATH" -type f -name "*.proto"  > "$SPEC_PATH/proto-files.txt"

echo "About to generate documentation for $SPEC"

protoc @"$SPEC_PATH/proto-files.txt" \
  --proto_path="$SPEC_PATH" \
  --proto_path="/protoc/include" \
  --proto_path="$PROTO_PATH" \
  --doc_out="$SPEC_PATH" --doc_opt=html,$GRPC_GEN_DOC_FILE

GCS_FILE_URL="https://storage.googleapis.com/$BUCKET_NAME/$SPEC/$GRPC_GEN_DOC_FILE"

echo "About to upload generated documentation for $SPEC to gs://$BUCKET_NAME"

curl --upload-file "$SPEC_PATH/$GRPC_GEN_DOC_FILE" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: text/html" \
  "$GCS_FILE_URL"

ARTIFACT_CONTENT=$(echo "$GCS_FILE_URL" | xxd -p | tr -d '\n')

echo "About to create grpc-doc-url artifact for $SPEC"
SPEC_WITHOUT_REV=$(echo "$SPEC" | sed 's/\@.*//')

registry rpc create-artifact \
  --parent="$SPEC_WITHOUT_REV" \
  --artifact_id="grpc-doc-url" \
  --artifact.mime_type="text/plain" \
  --artifact.contents="$ARTIFACT_CONTENT"

echo "Finished processing $SPEC"
