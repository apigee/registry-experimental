import grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2 as rs
import argparse
from metrics import vocabulary_pb2
from comparison import Comparison
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    word_group_pb2 as wg,
)
from google.cloud.apigeeregistry.v1 import registry_models_pb2 as rm
def main():
    # Creating registry client
    channel = grpc.insecure_channel("localhost:8080")
    stub = registry_service_pb2_grpc.RegistryStub(channel)

    parser = argparse.ArgumentParser()

    # add project_name and spec_name commandline arguments.
    parser.add_argument(
        "--spec_name",
        type=str,
        required=True,
        help="Name of the spec to compute comparison for",
    )
    parser.add_argument(
        "--project_name",
        type=str,
        required=True,
        help="Name of the project to extract WordGroups and compute comparison for",
    )

    args = parser.parse_args()

    spec_name = args.spec_name
    project_name = args.project_name
    print(spec_name, project_name)

    # Get vocabulary artifacts to generate new spec words
    try:
        response = stub.ListArtifacts(
            rs.ListArtifactsRequest(
                parent=spec_name,
                filter="name.contains(\"vocabulary\")"
            ))
    except grpc.RpcError as rpc_error:
        print(f"Received RPC error: code={rpc_error.code()} message={rpc_error.details()}")
    
    new_words = []
    for artifact in response.artifacts:
        contents = stub.GetArtifactContents(
            rs.GetArtifactContentsRequest(
                name=artifact.name,
            )
        )
        vocab = vocabulary_pb2.Vocabulary()
        vocab.ParseFromString(contents.data)
        for entry in vocab.schemas:
            if type(entry.word) == str and "." not in entry.word and len(entry.word) > 2:
                new_words.append(entry.word)
        for entry in vocab.properties:
            if type(entry.word) == str and "." not in entry.word and len(entry.word) > 2:
                new_words.append(entry.word)
        for entry in vocab.operations:
            if type(entry.word) == str and "." not in entry.word and len(entry.word) > 2:
                new_words.append(entry.word)
        for entry in vocab.parameters:
            if type(entry.word) == str and "." not in entry.word and len(entry.word) > 2:
                new_words.append(entry.word)

    # get wordgroups and noise_words to compare against and generate a report. 
    try:
        response = stub.ListArtifacts(
            rs.ListArtifactsRequest(
                parent="projects/"
                + project_name
                + "/locations/global/apis/-/versions/-/specs/-",
                filter='name.contains("ConsistencyReport")',
            )
        )
    except grpc.RpcError as rpc_error:
        print(
            f"Failed to fetch WordGroup artifacts, RPC error: code={rpc_error.code()} message={rpc_error.details()}"
        )
        return None

    word_groups = []
    noise_words = None
    for artifact in response.artifacts:
        contents = stub.GetArtifactContents(
            rs.GetArtifactContentsRequest(name=artifact.name)
        )

        wrdgrp = wg.WordGroup()
        wrdgrp.ParseFromString(contents.data)
        if wrdgrp.id == "NOISE_WORDS":
            noise_words = wrdgrp
            continue
        word_groups.append(wrdgrp)

    # call the comparison class
    cmprsn = Comparison(stub = stub, new_words=new_words, word_groups=word_groups, noise_words=noise_words)

    # generate a consistency report
    report = cmprsn.generate_consistency_report()
    print(report)


    ## upload the report
    artifact = rm.Artifact(
        name="projects/"
        + project_name
        + "/locations/global/artifacts/consistency-report",
        mime_type="application/octet-stream;type=google.cloud.apigeeregistry.applications.v1alpha1.consistency.WordGroup",
        contents=report.SerializeToString(),
    )

    createArtifactRequest = rs.CreateArtifactRequest(
        parent="projects/" + project_name + "/locations/global",
        artifact=artifact,
        artifact_id=report.id,
    )
    try:
        stub.CreateArtifact(createArtifactRequest)

    except grpc.RpcError as rpc_error:
        err = rpc_error.code()
        if err != grpc.StatusCode.ALREADY_EXISTS:
            print(f"Received RPC error: code= {err} message= {rpc_error.details()}")

        else:
            replaceArtifactRequest = rs.ReplaceArtifactRequest(artifact=artifact)

            try:
                stub.ReplaceArtifact(replaceArtifactRequest)

            except grpc.RpcError as rpc_error:
                err = rpc_error.code()
                print(
                    f"Received RPC error: code= {err} message= {rpc_error.details()}"
                )

if __name__ == "__main__":
    main()
