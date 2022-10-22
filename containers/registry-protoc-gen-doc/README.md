#Generating Documentation from protos

We will be creating a custom controller which can generate html documentation
using [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and store
the generated markup to an artifact `grpc-documentation` on the spec object.

1. Create a Google storage bucket to use for this setup.
2. Update the manifest file `registry-protoc-gen-doc-controller.yaml` with the name of the bucket (replace grpc-docs).
3. Upload the manifest for this controller to the Registry project
   ```
   PROJECT_ID=<GCP_PROJECT_NAME>
   registry upload manifest ./registry-protoc-gen-doc-manifest.yaml --project-id=${PROJECT_ID}
   ```
4. Deploy to controller to GKE 
   1. Create a GKE cluster with [WorkLoad Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) enabled.
   2. Create a IAM Service account
      ```shell
    
      gcloud iam service-accounts create registry-admin-gsa --project=${PROJECT_ID}
      
      gcloud projects add-iam-policy-binding $PROJECT_ID \
      --member "serviceAccount:registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/apigeeregistry.admin"
      
      gcloud projects add-iam-policy-binding $PROJECT_ID \
      --member "serviceAccount:registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/storage.objectAdmin"
      ```
   3. Configure workload identity for the service account
         ```
         gcloud iam service-accounts add-iam-policy-binding registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com \
         --role roles/iam.workloadIdentityUser \
         --member "serviceAccount:${PROJECT_ID}.svc.id.goog[registry-custom/registry-admin-ksa]"
      
         ``` 
   4. Apply the kubernetes configuration to the cluster
      ```
        kubectl apply -f registry-protoc-gen-doc-controller.yaml
      ```
   5. Annotate the Kubernetes service account with the email address of the IAM service account.
      ```
        kubectl annotate serviceaccount registry-admin-ksa \
          --namespace registry-custom \
         iam.gke.io/gcp-service-account=registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com
      ```

## Delete all the GRPC documentation artifacts
```shell
  #Delete grpc-doc-url artifact for all specs in registry 
  registry delete apis/-/versions/-/specs/-/artifacts/grpc-doc-url
  #Delete the files from the GCS bucket
  gsutil rm -rf gs://grpc-docs/*
```