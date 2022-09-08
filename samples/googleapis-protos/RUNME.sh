#!/bin/sh

PROJECT=googleapis

registry rpc admin delete-project --name projects/$PROJECT --force
registry rpc admin create-project --project_id $PROJECT
registry upload bulk protos ~/googleapis --project-id $PROJECT
registry apply -f artifacts -R --parent projects/$PROJECT/locations/global
registry compute conformance projects/$PROJECT/locations/global/apis/-/versions/-/specs/-
registry compute score projects/$PROJECT/locations/global/apis/-/versions/-/specs/-
registry compute scorecard projects/$PROJECT/locations/global/apis/-/versions/-/specs/-

