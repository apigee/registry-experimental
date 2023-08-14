#!/bin/bash
#
# Copyright 2022 Google LLC.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Create a fresh project for the Google API protos.
PROJECT=googleapis
registry rpc admin delete-project --name projects/$PROJECT --force
registry rpc admin create-project --project_id $PROJECT

# This upload assumes that github.com/googleapis/googlapis is checked out at the path below.
registry upload bulk protos ~/googleapis --project-id $PROJECT
# Apply artifacts that configure the style guide and scoring.
registry apply -f artifacts -R --parent projects/$PROJECT/locations/global
# Compute conformance reports and associated scores.
registry compute conformance projects/$PROJECT/locations/global/apis/-/versions/-/specs/-
registry compute score projects/$PROJECT/locations/global/apis/-/versions/-/specs/-
registry compute scorecard projects/$PROJECT/locations/global/apis/-/versions/-/specs/-