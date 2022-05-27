import grpc
from google.cloud.apigeeregistry.v1 import admin_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import admin_service_pb2
from google.cloud.apigeeregistry.v1 import admin_models_pb2
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from google.cloud.apigeeregistry.v1 import registry_models_pb2
#import admin_service_pb2_grpc

def main():
    # Creating admin client
    channel = grpc.insecure_channel('localhost:8080')
    stub = admin_service_pb2_grpc.AdminStub(channel)

    # Creating a project
    try:
        response = stub.CreateProject(
            admin_service_pb2.CreateProjectRequest(
                project=admin_models_pb2.Project(name="demo"),
                project_id="demo"))
        print(f"CreateProject response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")
    # Get a project
    try:
        response = stub.GetProject(
            admin_service_pb2.GetProjectRequest(name="projects/demo"))
        print(f"GetProject response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")


    # Creating registry client
    channel = grpc.insecure_channel('localhost:8080')
    stub = registry_service_pb2_grpc.RegistryStub(channel)
    # Creating an API
    try:
        response = stub.CreateApi(
            registry_service_pb2.CreateApiRequest(
                parent="projects/demo/locations/global",
                api=registry_models_pb2.Api(
                    name="projects/demo/locations/global/apis/petstore",
                    description="demo API petstore",
                ),
                api_id="petstore"))
        print(f"CreateApi response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")

    # Get an API
    try:
        response = stub.GetApi(
            registry_service_pb2.GetApiRequest(
                name="projects/demo/locations/global/apis/petstore",
            ))
        print(f"GetApi response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")

    # Creating a Version
    try:
        response = stub.CreateApiVersion(
            registry_service_pb2.CreateApiVersionRequest(
                parent="projects/demo/locations/global/apis/petstore",
                api_version=registry_models_pb2.ApiVersion(
                    name="projects/demo/locations/global/apis/petstore/versions/1.0.0",
                    description="demo API petstore v 1.0.0",
                ),
                api_version_id="1.0.0"))
        print(f"CreateApiVersion response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")

    # Get an ApiVersion
    try:
        response = stub.GetApiVersion(
            registry_service_pb2.GetApiVersionRequest(
                name="projects/demo/locations/global/apis/petstore/versions/1.0.0",
            ))
        print(f"GetApiVersion response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")
    
    # Creating an ApiSpec
    try:
        response = stub.CreateApiSpec(
            registry_service_pb2.CreateApiSpecRequest(
                parent="projects/demo/locations/global/apis/petstore/versions/1.0.0",
                api_spec=registry_models_pb2.ApiSpec(
                    name="projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
                    description="openapi.yaml",
                    mime_type="application/x.openapi+gzip;version=3",
                ),
                api_spec_id="openapi.yaml"))
        print(f"CreateApiSpec response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")
    
    # Get an ApiSpec
    try:
        response = stub.GetApiSpec(
            registry_service_pb2.GetApiSpecRequest(
                name="projects/demo/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
            ))
        print(f"GetApiSpec response: {response}")
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")

if __name__ == "__main__":
    main()