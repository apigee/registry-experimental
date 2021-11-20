lite:
	go install ./...

all:
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-APX.sh
	go install ./...

apg:
	./tools/GENERATE-APX.sh
	go install ./cmd/apx

protos:
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-ENVOY-DESCRIPTORS.sh

test:
	go test ./...

build-workers:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
	gcloud config set project ${REGISTRY_PROJECT_IDENTIFIER}
	gcloud builds submit --config deployments/worker-setup/cloudbuild.yaml \
    --substitutions _REGISTRY_PROJECT_IDENTIFIER="${REGISTRY_PROJECT_IDENTIFIER}"

deploy-workers:
	./deployments/worker-setup/DEPLOY-WORKERS.sh
