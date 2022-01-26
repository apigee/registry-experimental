# Registry Spec renderer
This container allows for rendering of specs from the API Registry.

### To run this against a self hosted version of Apigee Registry you will need to 
define 'APG_REGISTRY_ADDRESS' environment variable
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

### Running this service against hosted API Registry service 
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