import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2


class ExtractWords:
    def __init__(self, stub, project_name):
        self.stub = stub
        self.project_name = project_name

    def extract_vocabs(self):
        stub = self.stub
        try:
            response = stub.ListArtifacts(
                registry_service_pb2.ListArtifactsRequest(
                    parent="projects/" + self.project_name + "/locations/global/apis/-/versions/-/specs/-",
                    filter="name.contains(\"vocabulary\")"
                ))
        except grpc.RpcError as rpc_error:
            print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")

        vocabs = []
        for artifact in response.artifacts:
            contents = stub.GetArtifactContents(
                registry_service_pb2.GetArtifactContentsRequest(
                    name=artifact.name,
                )
            )
            vocab = vocabulary_pb2.Vocabulary()
            
            try:
                vocabs.append(vocab.ParseFromString(contents.data))
            except Exception as e:
                print(e, " Getting vocabulary for ", artifact.name, "failed!")
                continue

        if len(vocabs) < 1:
            return None

        return vocabs

    def get_vocabs(self):
        vocabs = self.extract_vocabs()
        if vocabs is None:
            return None

        words = []
        for vocab in vocabs:
            for entry in vocab.schemas:
                for _ in range(entry.count):
                    words.append(entry.word)
            for entry in vocab.properties:
                for _ in range(entry.count):
                    words.append(entry.word)
            for entry in vocab.operations:
                for _ in range(entry.count):
                    words.append(entry.word)
            for entry in vocab.parameters:
                for _ in range(entry.count):
                    words.append(entry.word)
        return words
