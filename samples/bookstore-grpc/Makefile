
PROJECT := $(shell gcloud config get project)

lite:
	go install ./...

all:
	mkdir -p ./cmd/bookstore ./rpc ./gapic
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-CLI.sh
	go get ./...
	go install ./...

clean:
	rm -rf cmd/bookstore/*.go gapic/*.go rpc/*.go third_party proto.pb docs

cloudrun:
	gcloud run deploy bookstore --source .

apigateway:
	./tools/GENERATE-DESCRIPTOR.sh
	./tools/GENERATE-APICONFIG.sh
	gcloud api-gateway api-configs create bookstore --api=bookstore --project=$(PROJECT) --grpc-files=proto.pb,api_config.yaml
	gcloud api-gateway gateways create bookstore --api=bookstore --api-config=bookstore --location=us-west2 --project=$(PROJECT)

test:
	go clean -testcache
	go test . -v
