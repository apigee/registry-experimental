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

# deploy registry-server on CloudRun
deploy:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
	gcloud run deploy registry-backend --image gcr.io/${REGISTRY_PROJECT_IDENTIFIER}/registry-server:latest --platform managed

# deploy registry-server on GKE
deploy-gke:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
ifeq ($(LB),internal)
	./deployments/registry-server/gke/DEPLOY-TO-GKE.sh deployments/registry-server/gke/service-internal.yaml
else
	./deployments/gke/DEPLOY-TO-GKE.sh
endif

# Actions for controller
deploy-controller-job:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
ifndef REGISTRY_MANIFEST_ID
	@echo "Error! REGISTRY_MANIFEST_ID must be set."; exit 1
endif
	gcloud container clusters get-credentials registry-backend --zone us-central1-a
	envsubst < deployments/controller/gke-job/cron-job.yaml | kubectl apply -f -

deploy-controller-dashboard:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
	./deployments/controller/dashboard/DEPLOY.sh

deploy-controller: deploy-controller-job deploy-controller-dashboard

build-workers:
ifndef REGISTRY_PROJECT_IDENTIFIER
	@echo "Error! REGISTRY_PROJECT_IDENTIFIER must be set."; exit 1
endif
	gcloud config set project ${REGISTRY_PROJECT_IDENTIFIER}
	gcloud builds submit --config deployments/worker-setup/cloudbuild.yaml \
    --substitutions _REGISTRY_PROJECT_IDENTIFIER="${REGISTRY_PROJECT_IDENTIFIER}"

deploy-workers:
	./deployments/worker-setup/DEPLOY-WORKERS.sh
