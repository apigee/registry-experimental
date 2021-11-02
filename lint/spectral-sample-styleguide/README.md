## Summary

The following repo contains a sample Spectral style guide that governs violations on OpenAPI Specs, along witha sample OpenAPI spec which violates multiple rules. This example is used to illustrate the functionality of style guides in the API registry as well as the generation of conformance reports. The test-description of this project teaches how to upload the style guide to the registry, upload the manifest to the registry, compute the conformance report, and the conformance report.

## Test

- Have registry running on `localhost:8080`.
- Ensure that both `petstore.yaml` and `styleguide.yaml` from this directory are located in `~`.
- Compile the registry project by running `make` in the root directory

```bash
apg registry delete-project --name projects/example

apg registry create-project --project_id example --insecure --address localhost:8080

apg registry create-api --api_id petstore --parent=projects/example/locations/global

apg registry create-api-version --api_version_id v1 --parent=projects/example/locations/global/apis/petstore

registry upload spec --version=projects/example/locations/global/apis/petstore/versions/v1 --style=openapi ~/petstore.yaml

registry upload styleguide ~/styleguide.yaml --project-id=example

registry compute conformance projects/example/locations/global/apis/petstore/versions/v1/specs/petstore.yaml --plugin=true

registry get projects/example/locations/global/apis/petstore/versions/v1/specs/petstore.yaml/artifacts/conformance-apilinterstyleguide --contents
```
