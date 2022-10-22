#Generating Documentation from protos

We will be creating a custom controller which can generate html documentation
using [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and store
the generated markup to an artifact `grpc-documentation` on the spec object.

1. We have to generate a manifest for this controller
```

PROJECT_ID=<GCP_PROJECT_NAME>
registry upload manifest ./registry-grpc-doc-gen-manifest.yaml --project-id=${PROJECT_ID}
```

2. Deploy to controller to GKE 
   1. Create a GKE cluster with [WorkLoad Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) enabled.
   2. Create a IAM Service account
      ```shell
     
      gcloud iam service-accounts create registry-admin-gsa --project=${PROJECT_ID}
      
      gcloud projects add-iam-policy-binding $PROJECT_ID \
      --member "serviceAccount:registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/apigeeregistry.admin"
      ```
   3. Configure workload identity for the service account
         ```
         gcloud iam service-accounts add-iam-policy-binding registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com \
         --role roles/iam.workloadIdentityUser \
         --member "serviceAccount:${PROJECT_ID}.svc.id.goog[registry-custom/registry-admin-ksa]"
      
         ``` 
   4. Apply the kubernetes configuration to the cluster
      ```
        kubectl apply -f registry-grpc-doc-gen-controller.yaml
      ```
   5. Annotate the Kubernetes service account with the email address of the IAM service account.
      ```
        kubectl annotate serviceaccount registry-admin-ksa \
          --namespace registry-custom \
         iam.gke.io/gcp-service-account=registry-admin-gsa@${PROJECT_ID}.iam.gserviceaccount.com
      ```

