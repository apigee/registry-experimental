# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from google.cloud.apigeeregistry.v1 import admin_models_pb2 as google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2
from google.cloud.apigeeregistry.v1 import admin_service_pb2 as google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2
from google.longrunning import operations_pb2 as google_dot_longrunning_dot_operations__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


class AdminStub(object):
    """The Admin service supports setup and operation of an API registry.
    It is typically not included in hosted versions of the API.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetStatus = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/GetStatus',
                request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Status.FromString,
                )
        self.GetStorage = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/GetStorage',
                request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Storage.FromString,
                )
        self.MigrateDatabase = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/MigrateDatabase',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.MigrateDatabaseRequest.SerializeToString,
                response_deserializer=google_dot_longrunning_dot_operations__pb2.Operation.FromString,
                )
        self.ListProjects = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/ListProjects',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsRequest.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsResponse.FromString,
                )
        self.GetProject = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/GetProject',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.GetProjectRequest.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
                )
        self.CreateProject = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/CreateProject',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.CreateProjectRequest.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
                )
        self.UpdateProject = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/UpdateProject',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.UpdateProjectRequest.SerializeToString,
                response_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
                )
        self.DeleteProject = channel.unary_unary(
                '/google.cloud.apigeeregistry.v1.Admin/DeleteProject',
                request_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.DeleteProjectRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                )


class AdminServicer(object):
    """The Admin service supports setup and operation of an API registry.
    It is typically not included in hosted versions of the API.
    """

    def GetStatus(self, request, context):
        """GetStatus returns the status of the service.
        (-- api-linter: core::0131::request-message-name=disabled
        aip.dev/not-precedent: Not in the official API. --)
        (-- api-linter: core::0131::method-signature=disabled
        aip.dev/not-precedent: Not in the official API. --)
        (-- api-linter: core::0131::http-uri-name=disabled
        aip.dev/not-precedent: Not in the official API. --)
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetStorage(self, request, context):
        """GetStorage returns information about the storage used by the service.
        (-- api-linter: core::0131::request-message-name=disabled
        aip.dev/not-precedent: Not in the official API. --)
        (-- api-linter: core::0131::method-signature=disabled
        aip.dev/not-precedent: Not in the official API. --)
        (-- api-linter: core::0131::http-uri-name=disabled
        aip.dev/not-precedent: Not in the official API. --)
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def MigrateDatabase(self, request, context):
        """MigrateDatabase attempts to migrate the database to the current schema.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListProjects(self, request, context):
        """ListProjects returns matching projects.
        (-- api-linter: standard-methods=disabled --)
        (-- api-linter: core::0132::method-signature=disabled
        aip.dev/not-precedent: projects are top-level resources. --)
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetProject(self, request, context):
        """GetProject returns a specified project.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CreateProject(self, request, context):
        """CreateProject creates a specified project.
        (-- api-linter: standard-methods=disabled --)
        (-- api-linter: core::0133::http-uri-parent=disabled
        aip.dev/not-precedent: Project has an implicit parent. --)
        (-- api-linter: core::0133::method-signature=disabled
        aip.dev/not-precedent: Project has an implicit parent. --)
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UpdateProject(self, request, context):
        """UpdateProject can be used to modify a specified project.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DeleteProject(self, request, context):
        """DeleteProject removes a specified project and all of the resources that it
        owns.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AdminServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.GetStatus,
                    request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Status.SerializeToString,
            ),
            'GetStorage': grpc.unary_unary_rpc_method_handler(
                    servicer.GetStorage,
                    request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Storage.SerializeToString,
            ),
            'MigrateDatabase': grpc.unary_unary_rpc_method_handler(
                    servicer.MigrateDatabase,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.MigrateDatabaseRequest.FromString,
                    response_serializer=google_dot_longrunning_dot_operations__pb2.Operation.SerializeToString,
            ),
            'ListProjects': grpc.unary_unary_rpc_method_handler(
                    servicer.ListProjects,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsRequest.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsResponse.SerializeToString,
            ),
            'GetProject': grpc.unary_unary_rpc_method_handler(
                    servicer.GetProject,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.GetProjectRequest.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.SerializeToString,
            ),
            'CreateProject': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateProject,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.CreateProjectRequest.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.SerializeToString,
            ),
            'UpdateProject': grpc.unary_unary_rpc_method_handler(
                    servicer.UpdateProject,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.UpdateProjectRequest.FromString,
                    response_serializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.SerializeToString,
            ),
            'DeleteProject': grpc.unary_unary_rpc_method_handler(
                    servicer.DeleteProject,
                    request_deserializer=google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.DeleteProjectRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'google.cloud.apigeeregistry.v1.Admin', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Admin(object):
    """The Admin service supports setup and operation of an API registry.
    It is typically not included in hosted versions of the API.
    """

    @staticmethod
    def GetStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/GetStatus',
            google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Status.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetStorage(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/GetStorage',
            google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Storage.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def MigrateDatabase(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/MigrateDatabase',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.MigrateDatabaseRequest.SerializeToString,
            google_dot_longrunning_dot_operations__pb2.Operation.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListProjects(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/ListProjects',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsRequest.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.ListProjectsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetProject(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/GetProject',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.GetProjectRequest.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CreateProject(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/CreateProject',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.CreateProjectRequest.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UpdateProject(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/UpdateProject',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.UpdateProjectRequest.SerializeToString,
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2.Project.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def DeleteProject(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/google.cloud.apigeeregistry.v1.Admin/DeleteProject',
            google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__service__pb2.DeleteProjectRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
