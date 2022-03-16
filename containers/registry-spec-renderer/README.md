# Registry Spec renderer
This container allows for rendering of specs from the API Registry.

### To run this service on a GCE instance run the following command:
```
export REGISTRY_PROJECT_IDENTIFIER=$(gcloud config list --format 'value(core.project)')
gcloud iam service-accounts create registry-reader \
    --description="Registry Reader" \
    --display-name="Registry Reader"

gcloud projects add-iam-policy-binding $REGISTRY_PROJECT_IDENTIFIER \
    --member="serviceAccount:registry-reader@$REGISTRY_PROJECT_IDENTIFIER.iam.gserviceaccount.com" \
    --role="roles/apigeeregistry.viewer"

gcloud compute firewall-rules create registry-renderer-service-fw \
    --action allow \
    --target-tags registry-spec-renderer \
    --source-ranges 0.0.0.0/0 \
    --rules tcp:80


gcloud compute instances create-with-container registry-renderer-instance1 \
	--machine-type=e2-micro  --tags=registry-spec-renderer,http-server \
	--scopes=https://www.googleapis.com/auth/cloud-platform \
	--restart-on-failure --service-account=registry-reader@$REGISTRY_PROJECT_IDENTIFIER.iam.gserviceaccount.com\
    --container-image ghcr.io/apigee/registry-spec-renderer:latest
```

### To run this against the opensource version of Apigee Registry on GKE you will need to:  

1. Create a namespace for registry-spec-renderer
    ```
   kubectl create ns registry-spec-renderer
   ```
2. Store the registry service information to configmap
   ```
    kubectl create configmap registry-service-config -n registry-spec-renderer \
   --from-literal=APG_REGISTRY_ADDRESS=registry-service:8888
   ```
3. Apply the deployment file
   ```
    kubectl apply -f kubernetes/deployment-self-hosted.yaml -n registry-spec-renderer
   ```

### Running this service against hosted API Registry service on GKE: 

*Instead of using key you can configure workload identity*
1. you will need to create a service account with the 'roles/apigeeregistry.viewer' role

2. You can download the service account key and rename the file to service-account.json

3. Create a namespace for registry-spec-renderer
    ```
   kubectl create ns registry-spec-renderer
   ```
4. Store the service-account.json to secret 
 ```
   kubectl create secret generic registry-spec-renderer-sa-key \
   --from-file service-account.json -n registry-spec-renderer
   ```
5. Apply the deployment 
   ```
    kubectl apply -f kubernetes/deployment.yaml -n registry-spec-renderer
   ```
