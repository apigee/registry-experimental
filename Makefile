all:
	go install ./...

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
