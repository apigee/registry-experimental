#Generating Documentation from protos

We will be creating a custom controller which can generate html documentation
using [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and store
the generated markup to an artifact `grpc-documentation` on the spec object.

1. We have to generate a manifest for this controller
```
registry upload manifest ./registry-grpc-doc-gen-manifest.yaml --project-id=<GCP PROJECT>
```

2. 
