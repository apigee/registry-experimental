# Mock server for OpenAPI specs

This application will create a mock server which serves all the OpenAPI specs
stored in the Apigee Registry.

For a registry spec 'projects/openapi/locations/global/apis/benchtest-1/versions/v1/specs/spec-1'
The mock service will be available at
http://localhost:3000/projects/openapi/locations/global/apis/benchtest-1/versions/v1/specs/spec-1

Mock service is also available at http://localhost:3000/mock 
by passing the header specifying the spec path
`apigee-registry-spec: projects/openapi/locations/global/apis/benchtest-1/versions/v1/specs/spec-1`


[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)