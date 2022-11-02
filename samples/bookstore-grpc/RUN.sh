#!/bin/sh

# build locally
make all

# test locally
bookstore-server &
export BOOKSTORE_BOOKSTORE_ADDRESS=localhost:8080
export BOOKSTORE_BOOKSTORE_INSECURE=1
go test .

# upload protos for the API into the registry
registry upload bulk protos . --project-id timburks-test

# deploy the backend on cloud run
gcloud run deploy bookstore --source .

# register the backend deployment
./tools/GENERATE-BACKEND-DEPLOYMENT.sh
registry apply -f backend-deployment.yaml

# test the backend deployment
export BOOKSTORE_BOOKSTORE_ADDRESS=$()
export BOOKSTORE_BOOKSTORE_INSECURE=1
go test .

# compute the file descriptor set (proto.pb) for the API
# (requires registry-experimental and protoc)
registry-experimental compute descriptor projects/timburks-test/locations/global/apis/bookstore/versions/v1/specs/examples-bookstore-v1.zip

# deploy a gateway that uses the backend
./tools/GENERATE-DESCRIPTOR.sh
./tools/GENERATE-APICONFIG.sh
gcloud api-gateway api-configs create bookstore --api=bookstore --project=$(PROJECT) --grpc-files=proto.pb,api_config.yaml
gcloud api-gateway gateways create bookstore --api=bookstore --api-config=bookstore --location=us-west2 --project=$(PROJECT)

# register the gateway deployment
./tools/GENERATE-GATEWAY-DEPLOYMENT.sh
registry apply -f gateway-deployment.yaml

# test the gateway deployment
export BOOKSTORE_BOOKSTORE_ADDRESS=$()
export BOOKSTORE_BOOKSTORE_INSECURE=1
go test .
