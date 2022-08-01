import grpc
import warnings
import os
import pandas as pd
import argparse
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2 as rs
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    consistency_report_pb2 as cr,
)


def generate_csv():

    channel = grpc.insecure_channel("localhost:8080")
    stub = registry_service_pb2_grpc.RegistryStub(channel)

    parser = argparse.ArgumentParser()

    # add project_name
    parser.add_argument(
        "--project_name",
        type=str,
        required=True,
        help="Name of project to the compute csv file for",
    )

    # add path to folder
    parser.add_argument(
        "--path",
        type=str,
        required=True,
        help="Path to the directory to save the CSV file",
    )

    # add file name
    parser.add_argument(
        "--csv_name", type=str, required=True, help="Name of the generated CSV file"
    )

    args = parser.parse_args()
    project_name = args.project_name
    path = args.path
    name = args.csv_name

    # get wordgroups and noise_words to compare against and generate a report.
    try:
        response = stub.ListArtifacts(
            rs.ListArtifactsRequest(
                parent="projects/"
                + project_name
                + "/locations/global/apis/-/versions/-/specs/-",
                filter='mime_type.contains("ConsistencyReport")',
            )
        )
    except grpc.RpcError as rpc_error:
        print(
            f"Failed to fetch Consistency Reports artifacts, RPC error: code={rpc_error.code()} message={rpc_error.details()}"
        )
        return None
    consistency_reports = []

    for artifact in response.artifacts:
        contents = stub.GetArtifactContents(
            rs.GetArtifactContentsRequest(name=artifact.name)
        )
        report = cr.ConsistencyReport()
        report.ParseFromString(contents.data)
        consistency_reports.append(report)

    df = pd.DataFrame(columns=["new_word", "closest_cluster_id", "cluster_words"])
    csv_file_name = name + ".csv"
    path = os.path.join(path, csv_file_name)
    for i in range(len(consistency_reports)):
        report = consistency_reports[i]
        report.current_variations.sort(key=lambda x: x.term)
        for variation in report.current_variations:
            df.loc[len(df)] = [
                variation.term,
                variation.cluster.id,
                list(variation.cluster.word_frequency.keys()),
            ]

    df.to_csv(path)


if __name__ == "__main__":
    generate_csv()
