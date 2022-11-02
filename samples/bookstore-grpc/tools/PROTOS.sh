#!/bin/bash

SERVICE_PROTOS=(
	examples/bookstore/v1/*.proto
)

COMMON_PROTOS_PATH='third_party/api-common-protos'
if [ ! -d $COMMON_PROTOS_PATH ]
then
	git clone https://github.com/googleapis/api-common-protos $COMMON_PROTOS_PATH
fi
