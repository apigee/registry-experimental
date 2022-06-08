set -e
# clone registry protos
REGISTRY_PROTOS_PATH='registry-protos'
if [ ! -d $REGISTRY_PROTOS_PATH ]
then
    git clone https://github.com/apigee/registry $REGISTRY_PROTOS_PATH
fi
# clone common protos
COMMON_PROTOS_PATH='third_party/api-common-protos'
if [ ! -d $COMMON_PROTOS_PATH ]
then
    git clone https://github.com/googleapis/api-common-protos $COMMON_PROTOS_PATH
fi
PYTHON_TYPES_PATH='consistency/rpc'
ALL_PROTOS=$REGISTRY_PROTOS_PATH/google/cloud/apigeeregistry/v1/*.proto
for proto in ${ALL_PROTOS[@]}; do
	echo "Generating Python types for $proto"
	protoc $proto --proto_path=$REGISTRY_PROTOS_PATH --proto_path=$COMMON_PROTOS_PATH --python_out=$PYTHON_TYPES_PATH
done