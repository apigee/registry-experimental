import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2


class ExtractWords:
    def __init__(self, stub):
        self.stub = stub

    def extract_vocabs(self):
        stub = self.stub
            # Get vocabulary artifacts
        vocabs = []
        try:
            response = stub.ListArtifacts(
                registry_service_pb2.ListArtifactsRequest(
                    parent="projects/-/locations/global/apis/-/versions/-/specs/-",
                    filter='name.contains("vocabulary")',
                )
            )
        except grpc.RpcError as rpc_error:
            print(
                f"Failed to fetch vocabulary artifacts, RPC error: code={rpc_error.code()} message={rpc_error.details()}"
            )
            return None

        vocabs = []
        for artifact in response.artifacts:
            contents = stub.GetArtifactContents(
                registry_service_pb2.GetArtifactContentsRequest(name=artifact.name)
            )

            vocab = vocabulary_pb2.Vocabulary()
            vocab.ParseFromString(contents.data)
            vocabs.append(vocab)
            
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
