import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_models_pb2 as rm
from google.cloud.apigeeregistry.v1 import registry_service_pb2 as rs 
from word_extraction import ExtractWords
from clustering import ClusterWords 
from datetime import datetime 
import argparse

def main ():
    # Creating registry client
    channel = grpc.insecure_channel('localhost:8080')
    stub = registry_service_pb2_grpc.RegistryStub(channel)

    parser = argparse.ArgumentParser()

    parser.add_argument("--project_name", type=str, required=True, help= 'Name of the project to compute clusters for')

    args = parser.parse_args()
    # type cast the arguments as default is string type
    project_name = args.project_name

    extrct = ExtractWords(stub=stub, project_name=project_name)
    try:
            start_time = datetime.now() 

            words = extrct.get_vocabs()

            time_elapsed = datetime.now() - start_time 

            print('Time elapsed get vocabs (hh:mm:ss.ms) {}'.format(time_elapsed))
            
    except Exception as e:
                print(e, " \n Getting words failed")

    # clustr = ClusterWords(stub = stub, words=words)
    # try:
    #         start_time = datetime.now() 

    #         word_groups = clustr.create_word_groups()
    #         #print(word_groups)

    #         time_elapsed = datetime.now() - start_time 

    #         print('Time elapsed cluster words (hh:mm:ss.ms) {}'.format(time_elapsed))
        
        
    # except Exception as e:
    #             print(e, " \n Clustering words failed")
        
    #     ##### 
    #     # upload the wordGroups to the server 
    
    # for word_group in word_groups:
        
    #         id = "".join(filter(str.isalnum, word_group.id.lower()))

    #         artifact = rm.Artifact(
    #         name = "projects/" + project_name + "/locations/global/artifacts/wordgroup-" + id,
    #         mime_type = "application/octet-stream;type=google.cloud.apigeeregistry.applications.v1alpha1.consistency.WordGroup",
    #         contents = word_group.SerializeToString()
                
    #         )
            
    #         createArtifactRequest = rs.CreateArtifactRequest(
    #         parent = "projects/" + project_name + "/locations/global",
    #         artifact = artifact,
    #         artifact_id = "wordgroup-" + id
    #         )
    #         try:
    #             stub.CreateArtifact(
    #                 createArtifactRequest
    #               )
                
    #         except grpc.RpcError as rpc_error:
    #             err = rpc_error.code()
    #             if err != grpc.StatusCode.ALREADY_EXISTS:
    #                 print(f"Received RPC error: code= {err} message= {rpc_error.details()}")

    #             else:
    #                 replaceArtifactRequest = rs.ReplaceArtifactRequest (
    #                     artifact = artifact
    #                 )
    
    #                 try:
    #                     stub.ReplaceArtifact(
    #                         replaceArtifactRequest
    #                 )

    #                 except grpc.RpcError as rpc_error:
    #                     err = rpc_error.code()
    #                     print(f"Received RPC error: code= {err} message= {rpc_error.details()}")

if __name__ == '__main__':
    main()
