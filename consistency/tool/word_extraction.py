from urllib import response
import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2
class ExtractWords:
    def __init__(self, stub, linearize=True):
      self.stub = stub
      self.linearize = linearize
    def extract_vocabs(self):
      stub = self.stub
      try:
          response = stub.ListArtifacts(
              registry_service_pb2.ListArtifactsRequest(
                  parent="projects/-/locations/global/apis/-/versions/-/specs/-",
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
          vocabs.append(vocab.ParseFromString(contents.data))

      return vocabs
    def get_vocabs(self):
        vocabs = self.extract_vocabs(self)

        words = [] if self.linearize else {}
          
        if self.linearize:
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
 
          # Put words in a map of frequencies.  
        else:
            for vocab in vocabs:
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
 

