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

FROM pseudomuto/protoc-gen-doc:1.5 as protoc-gen-doc

FROM ghcr.io/apigee/registry-tools:latest

RUN apk update && apk add --no-cache protobuf \
    && mkdir -p /protoc/include/

COPY --from=protoc-gen-doc /usr/bin/protoc-gen-doc /bin/
COPY --from=protoc-gen-doc /usr/include/google /protoc/include/google/

RUN find /protoc/include/ -type f ! -iname "*.proto" -delete

COPY ./doc-gen.sh /

RUN chmod +x /doc-gen.sh

RUN git clone https://github.com/googleapis/api-common-protos /googleapis-common-protos
