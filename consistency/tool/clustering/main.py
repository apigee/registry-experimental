import grpc
from google.cloud.apigeeregistry.v1 import admin_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from word_extraction import ExtractWords
from clustering import ClusterWords
def main ():
   # Creating registry client
   channel = grpc.insecure_channel('localhost:8080')
   stub = registry_service_pb2_grpc.RegistryStub(channel)
   extrct = ExtractWords(stub=stub)
   try:
       words = extrct.get_vocabs()
   except Exception as e:
               print(e, " \n Getting words failed")
 
   clustr = ClusterWords(stub = stub, words=words)
   try:
       word_grups = clustr.create_word_groups()
       print(word_grups)
   except Exception as e:
               print(e, " \n Clustering words failed")
 
 
if __name__ == '__main__':
 main()
