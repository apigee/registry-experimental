#!/bin/bash
set -e

# clone registry protos
REGISTRY_PROTOS_PATH='consistency/registry-protos'
if [ ! -d $REGISTRY_PROTOS_PATH ]
then
    git clone https://github.com/apigee/registry $REGISTRY_PROTOS_PATH
fi

# clone vocabulary protos
GNOSTIC_PROTOS_PATH='consistency/gnostic-protos'
if [ ! -d $GNOSTIC_PROTOS_PATH ]
then
    git clone https://github.com/google/gnostic $GNOSTIC_PROTOS_PATH
fi

# clone common protos
COMMON_PROTOS_PATH='third_party/api-common-protos'
if [ ! -d $COMMON_PROTOS_PATH ]
then
    git clone https://github.com/googleapis/api-common-protos $COMMON_PROTOS_PATH
fi

ALL_PROTOS=(
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/applications/v1alpha1/*.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/internal/v1/*.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/*.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/controller/*.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/apihub/*.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/scoring/*.proto
	$GNOSTIC_PROTOS_PATH/metrics/*.proto
	google/cloud/apigeeregistry/v1/*.proto

)
SERVICE_PROTOS=(
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/registry_models.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/registry_service.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/admin_models.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/admin_service.proto
	$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/provisioning_service.proto
)

PYTHON_TYPES_PATH='consistency/rpc'
# # Generating grpc clients 
pip install grpcio-tools
echo "Generating gRPC client for ${ALL_PROTOS[*]}"
python -m grpc_tools.protoc --proto_path=$REGISTRY_PROTOS_PATH --proto_path=. --proto_path=$COMMON_PROTOS_PATH --proto_path=$GNOSTIC_PROTOS_PATH --python_out=$PYTHON_TYPES_PATH --grpc_python_out=$PYTHON_TYPES_PATH ${ALL_PROTOS[*]}

pip install consistency/rpc