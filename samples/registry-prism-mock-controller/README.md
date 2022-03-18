# Prism Mock controller

```
    MOCKSERVICE_DOMAIN=35-192-38-57.sslip.io
    REGISTRY_PROJECT_NAME=openapi
    CLOUDSDK_CONTAINER_CLUSTER=registry-mock-server-cluster
    CLOUDSDK_COMPUTE_ZONE=us-central1-c 
    envsubst < 02-prism-controller.yaml | kubectl apply -f -
```
