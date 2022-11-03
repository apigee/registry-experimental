#!/bin/bash
PROJECT=examples

echo Create a fresh project.
registry rpc admin delete-project --name projects/$PROJECT --force >& /dev/null
registry rpc admin create-project --project_id $PROJECT --json
registry config set registry.project $PROJECT

echo Create an entry for the Petstore API.
registry apply -f petstore.yaml

echo Upload each revision of the Petstore API spec.
registry apply -f petstore-openapi-r1.yaml
registry apply -f petstore-openapi-r2.yaml
registry apply -f petstore-openapi-r3.yaml

echo List all specs.
registry list apis/-/versions/-/specs/

echo List all spec revisions.
registry list apis/-/versions/-/specs/-@-

echo Compute complexity of the latest spec revisions.
registry compute complexity apis/-/versions/-/specs/-
registry list apis/-/versions/-/specs/-/artifacts
registry list apis/-/versions/-/specs/-@-/artifacts

echo Compute complexity of all spec revisions.
registry compute complexity apis/-/versions/-/specs/-@-
registry list apis/-/versions/-/specs/-/artifacts
registry list apis/-/versions/-/specs/-@-/artifacts

echo View the complexity of each revision.
for spec_revision in `registry list apis/petstore/versions/v1/specs/-@-`
do
        echo "$spec_revision"
        registry get $spec_revision/artifacts/complexity --print
done