# artifact-server

This small HTTP server calls a registry server backend to get the contents of
artifacts, which it then serves as JSON, YAML, or text (proto) documents
depending on the Accept type specified in request headers.

`registry get` remains the preferred method for getting artifact contents; this
adds the ability to get contents from a browser.

artifact-server reads the same configuration as the registry tool.
```
$ artifact-server &

$ curl -s http://localhost:8080/projects/sample/locations/global/artifacts/apihub-styleguide -H "Accept: text/json" | head -7
{
  "id": "apihub-styleguide",
  "kind": "StyleGuide",
  "mimeTypes": [
    "application/x.protobuf+zip"
  ],
  "guidelines": [

$ curl -s http://localhost:8080/projects/sample/locations/global/artifacts/apihub-styleguide -H "Accept: text/yaml" | head -7
id: apihub-styleguide
kind: StyleGuide
mimeTypes:
    - application/x.protobuf+zip
guidelines:
    - id: aip122
      rules:

$ curl -s http://localhost:8080/projects/sample/locations/global/artifacts/apihub-styleguide -H "Accept: text/plain" | head -7
id: "apihub-styleguide"
kind: "StyleGuide"
mime_types: "application/x.protobuf+zip"
guidelines: {
  id: "aip122"
  rules: {
    id: "camel-case-uris"
```