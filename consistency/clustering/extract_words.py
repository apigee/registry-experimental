import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2

class ExtractWords:
    def __init__(self, stub, linearize=True):
        self.stub = stub
        self.linearize = linearize

    def extract_words(self):

        stub = self.stub
        linearize = self.linearize

        try:
            response = stub.ListArtifacts(
                registry_service_pb2.ListArtifactsRequest(
                    parent="projects/-/locations/global/apis/-/versions/-/specs/-",
                    filter="name.contains(\"vocabulary\")"
                ))
        except grpc.RpcError as rpc_error:
            print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")
        
        for artifact in response.artifacts:
            contents = stub.GetArtifactContents(
                registry_service_pb2.GetArtifactContentsRequest(
                    name=artifact.name,
                )
            )

            vocab = vocabulary_pb2.Vocabulary()
            vocab.ParseFromString(contents.data)

            if linearize:
                words = []
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
            else:
                words = {}
                for entry in vocab.schemas:
                    word = entry.word
                    count= entry.count
                    if word not in words:
                        words[word] = 0 
                    words[word] += count

                for entry in vocab.properties:
                    word = entry.word
                    count= entry.count
                    if word not in words:
                        words[word] = 0 
                    words[word] += count

                for entry in vocab.operations:
                    word = entry.word
                    count= entry.count
                    if word not in words:
                        words[word] = 0 
                    words[word] += count

                for entry in vocab.parameters:
                    word = entry.word
                    count= entry.count
                    if word not in words:
                        words[word] = 0 
                    words[word] += count
                return words


