# Deploying Registry Server and Registry Viewer on GKE

The kubernetes directory has all the configuration to setup a registry server
and registry viewer in you Kubernetes cluster.

This deployment runs a pod with postgres database. It can be easily replaced 
with Cloud SQL for PostgreSQL.

Steps to install Apigee Registry & Viewer in your GKE Cluster:
1. Copy the sample folder and create a new folder
   > cp -rf kubenetes/sample kubenetes/registry-demo

2. We need a static IP for this demo: 
    ```
       gcloud compute addresses create registry-app-static-ip \
       --global \
       --ip-version IPV4 
   ```

   Please note down the above IP addresses.

3. You may choose to use your own DNS entry or use the wildcard sslip.io domain
    > This is required so that we can generate Google Managed SSL Certs.
    > 
    > e.g. if "registry-app-static-ip" is 1.2.3.4 you can use 1-2-3-4.sslip.io

    > If you are using your own DNS entries create an A record :    
      registry-demo.example.com  points to 1.2.3.4

   
   Replace the 1 occurrence of `registry-app.example.com` in kubernetes/registry-demo/patch.yaml with `1-2-3-4.sslip.io`
   or the custom domain you reserved for setting this demo

4. If you choose to use Cloud SQL with PostgreSQL you can modify the 
   `REGISTRY_DATABASE_CONFIG` entry and replace the string with the details
   You will have to additionally delete the registry-database pods once the setup is complete.

5. Create a Google oAuth Client ID are [here](https://console.cloud.google.com/apis/credentials/oauthclient)
   
    **Select the correct GCP Project.**
    Use the following values for the form:   
    - Application Type : Web Application
    - Name : API Registry Viewer
    - Javascript Authorized origins : Custom domain of the registry viewer 
        (e.g. https://1-2-3-4.sslip.io or https://registry-app.example.com)
    - Authorized redirect URIs: Custom domain of the registry viewer 
        (e.g. https://1-2-3-4.sslip.io or https://registry-app.example.com)

    Replace the `GOOGLE_SIGNIN_CLIENTID` in kubernetes/registry-demo/patch.yaml with the Client ID that was generated.

5. Ensure you are connected to the correct GKE cluster :
    ```
      gcloud container clusters get-credentials **cluster-name** --region **region** \
       --project **GCP-Project**
    ```
6. Run the following command to create the deployment setup:
    ```
      kubectl create ns api-registry

      kubectl apply -k kubernetes/registry-demo
    ```

7. The solution should be ready in about 10-15 minutes. 
   Check the status of the following:
   - Check the status of certs using :
     ```
        gcloud compute ssl-certificates list --global
     ```
   - Deploy the ingress rules
   - Health checks to pass

8. You should be able to access the viewer using https://1-2-3-4.sslip.io 
   or your custom domain
   
9. Get the IP address for the registry GRPC service using the following command:
    ```
   export APG_REGISTRY_ADDRESS=$(kubectl get svc -n api-registry registry-server-external-lb  -o jsonpath="{.status.loadBalancer.ingress[0].ip}:{.spec.ports[0].port}")
   export APG_ADMIN_ADDRESS=$APG_REGISTRY_ADDRESS
   export APG_ADMIN_INSECURE=1
   export APG_REGISTRY_INSECURE=1
    ```

10. Now you can interact with the registry tools using
    ```
    apg admin create-project --project_id=project1
    apg registry create-api --api_id=api1 --parent=projects/project1/locations/global
    apg registry list-apis --parent=projects/project1/locations/global --json
    ```
