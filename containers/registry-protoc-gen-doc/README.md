# Generating Documentation from protos

We will be running the registry controller with a custom manifest which can :
- generate html documentation using [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc)
- store the generated markup to Google Storage bucket
- store the reference to the GCS object, as an artifact `grpc-documentation`, on the spec object.

## Run this solution on location machine
1. Create a Google storage bucket to use for this setup.
2. Install [protoc](https://grpc.io/docs/protoc-installation/)
3. Install [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc)
4. Download Google APIs common protos using 
   ```
     git clone https://github.com/googleapis/api-common-protos
   ```
5. Updates the action field in the manifest file `registry-protoc-gen-doc-controller.yaml`
   ```
   /tmp/doc-gen.sh $resource.spec grpc-docs  /path/to/api-common-protos
   ```
   1. First parameter is the absolute path to doc-gen.sh file 
   2. Second parameter to the doc-gen.sh file is the name of the GCS bucket
   3. Third parameter is the spec reference 
   4. Fourth parameter is the path to the protos from https://github.com/googleapis/api-common-protos
      1. This could be referenced by any other common protos from your project
6. Upload the manifest for this controller to the Registry project
   ```shell
   export TOKEN=$(gcloud auth print-access-token)
   registry upload manifest ./registry-protoc-gen-doc-manifest.yaml
   ```
7. To generate the HTML markup and the artifacts on the spec run the following
    ```shell
    registry resolve artifacts/registry-protoc-gen-doc
    ```
8. List the artifacts using the below command
    ```shell
   registry list apis/-/versions/-/specs/-/artifacts/grpc-doc-url
    ```

## Run this solution on GKE
1. Create a Google storage bucket to use for this setup.
2. Update the manifest file `registry-protoc-gen-doc-controller.yaml` with the
   name of the bucket (replace grpc-docs).
3. Upload the manifest for this controller to the Registry project
   ```
   PROJECT_ID=<GCP_PROJECT_NAME>
   registry upload manifest ./registry-protoc-gen-doc-manifest.yaml --project-id=${PROJECT_ID}
   ```
4. Deploy the controller to GKE
    1. Create a GKE cluster
       with [WorkLoad Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
       enabled.
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
    5. Annotate the Kubernetes service account with the email address of the IAM
       service account.
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