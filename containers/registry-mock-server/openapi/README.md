# Mock server for OpenAPI specs

This application will create a mock server which serves all the OpenAPI specs
stored in the Apigee Registry.

For a registry spec 'projects/openapi/locations/global/apis/benchtest-1/versions/v1/specs/spec-1'
The mock service will be available at
http://localhost:3000/projects/openapi/locations/global/apis/benchtest-1/versions/v1/specs/spec-1

### To run this service on a GCE instance run the following command:
```
export REGISTRY_PROJECT_IDENTIFIER=$(gcloud config list --format 'value(core.project)')
gcloud iam service-accounts create registry-viewer \
    --description="Registry Reader" \
    --display-name="Registry Reader"

gcloud projects add-iam-policy-binding $REGISTRY_PROJECT_IDENTIFIER \
    --member="serviceAccount:registry-viewer@$REGISTRY_PROJECT_IDENTIFIER.iam.gserviceaccount.com" \
    --role="roles/apigeeregistry.viewer"

gcloud compute firewall-rules create registry-mock-service-fw \
    --action allow \
    --target-tags registry-mock-service \
    --source-ranges 0.0.0.0/0 \
    --rules tcp:80


gcloud compute instances create-with-container registry-mock-server-instance \
	--machine-type=e2-micro  --tags=registry-mock-service,http-server \
	--scopes=https://www.googleapis.com/auth/cloud-platform \
	--restart-on-failure --service-account=registry-viewer@$REGISTRY_PROJECT_IDENTIFIER.iam.gserviceaccount.com\
	--zone=us-central1-a \
    --container-image ghcr.io/giteshk-org/registry-prism-mock-server:main
```