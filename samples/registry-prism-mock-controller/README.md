# Registy Mock Generator

* It provides the ability to generate mock endpoints using Prism.
* A custom controller (extends the base registry controller framework) checks 
  the state of the Registry and deploys the prism mock service for every 
  OpenAPI specification.
* An artifact (`**_mock-prism-endpoint_**`) is generated, once the mock
  service is deployed successfully.


### Installations steps

The setup uses a GKE Cluster with ASM (Anthos Service Mesh). 
To use Istio please modify the configurations files in [kubernetes](kubenertes) 
directory. 
1. Create a GKE cluster, named `registry-mock-server-cluster` with ASM and workload identity enabled 
2. Connect to the cluster
3. Apply the ASM configuration
   ```
    kubectl apply -f kubernetes/01-asm-configuration.yaml
   ```
4. Create a service account with Registry Admin role (for the controller to write back artifacts)
   ```
   export REGISTRY_PROJECT_NAME=openapi
   
   gcloud config set project $REGISTRY_PROJECT_NAME
   
   gcloud iam service-accounts create registry-admin \
    --project=${REGISTRY_PROJECT_NAME}
   
   gcloud projects add-iam-policy-binding registry-admin \
    --member "serviceAccount:registry-admin@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com" \
    --role "roles/apigeeregistry.admin"

   gcloud projects add-iam-policy-binding registry-admin \
    --member "serviceAccount:registry-admin@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com" \
    --role "roles/logging.logWriter"

   gcloud iam service-accounts add-iam-policy-binding registry-admin@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${REGISTRY_PROJECT_NAME}.svc.id.goog[prism/registry-admin-sa]"

   kubectl annotate serviceaccount registry-viewer-sa \
    --namespace prism \
    iam.gke.io/gcp-service-account=registry-viewer@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com
   ```
5. Create a service account with the Apigee Registry Viewer permission (to read specs from the registry)
   ```
   export REGISTRY_PROJECT_NAME=openapi

   gcloud iam service-accounts create registry-viewer \
    --project=${REGISTRY_PROJECT_NAME}
   
   gcloud projects add-iam-policy-binding registry-viewer \
    --member "serviceAccount:registry-viewer@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com" \
    --role "roles/apigeeregistry.viewer"

   gcloud projects add-iam-policy-binding registry-viewer \
    --member "serviceAccount:registry-viewer@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com" \
    --role "roles/logging.logWriter"

   gcloud iam service-accounts add-iam-policy-binding registry-viewer@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${REGISTRY_PROJECT_NAME}.svc.id.goog[prism/registry-viewer-sa]"

   kubectl annotate serviceaccount registry-admin-sa \
    --namespace prism \
    iam.gke.io/gcp-service-account=registry-admin@${REGISTRY_PROJECT_NAME}.iam.gserviceaccount.com

   ```
6. Get the external IP address of Istio Ingress Gateway
    ```
    echo $(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    ```
7. Create a custom domain that points to this ip address. You can use sslip.io 
domain for demonstration purposes. For the IP address 3.4.5.6 use 3-4-5-6.sslip.io 
8. Upload the manifest to the registry
    ```
    registry upload manifest ./prism-manifest.yaml --project-id=${REGISTRY_PROJECT_NAME}
    ```
9. Verify the environment variables and apply the prism controller deployment.
```
    export APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
    export APG_REGISTRY_INSECURE=0
    export MOCKSERVICE_DOMAIN=3-4-5-6.sslip.io

    export CLOUDSDK_CONTAINER_CLUSTER=registry-mock-server-cluster
    export CLOUDSDK_COMPUTE_ZONE=us-central1-c 
    envsubst < kubernetes/02-prism-controller.yaml | kubectl apply -f -
```
10. The controller will create an artifact `prism-mock-endpoint` on every spec of type OpenAPI.

