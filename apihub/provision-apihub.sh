#!/bin/bash
#
# Copyright 2022 Google LLC. All Rights Reserved.
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

# This script automates the "Get started with API hub" flow described
# at https://cloud.google.com/apigee/docs/api-hub/get-started-api-hub.
# It handles prerequisites and provisions an instance of the API hub
# which includes a hosted instance of the Apigee Registry API.
# For more about the Registry Provisioning API that is used in this
# script, see https://cloud.google.com/apigee/docs/api-hub/provision.

# Prerequisites:
# - A Google Cloud account (USER_EMAIL) with billing enabled.
# - A Google Cloud project (PROJECT_ID).
# - The gcloud tool (https://cloud.google.com/sdk/gcloud).
# - A chosen location (RUNTIME_LOCATION) for the instance. Currently "us-central1"
#   is recommended. Other supported locations are listed at
#     https://cloud.google.com/apigee/docs/api-hub/provision
# - A chosen CMEK (https://cloud.google.com/kms/docs/cmek) key ring name
#   (KEY_RING_NAME) and key name (KEY_NAME).

# Usage:
#  provision-apihub.sh -u USER_EMAIL -p PROJECT_ID -l RUNTIME_LOCATION -r KEY_RING_NAME -k KEY_NAME

unset project_id
unset location
unset user
unset key_ring_name
unset key_name

while getopts p:l:u:r:k: opt; do
    case $opt in
        p) project_id=$OPTARG ;;
        l) location=$OPTARG ;;
        u) user=$OPTARG ;;
        r) key_ring_name=$OPTARG ;;
        k) key_name=$OPTARG ;;
        *)
            echo 'Error in command line parsing' >&2
            exit 1
    esac
done

if [[ -z "$project_id" ]] || [[ -z "$location" ]] || [[ -z "$user" ]]; then
    echo 'Missing the required parameters: -p, -l or -u' >&2
    exit 1
fi

if [[ -z "$key_ring_name" ]]; then
    key_ring_name="apihub-key-ring"
fi

if [[ -z "$key_name" ]]; then
    key_name="apihub-key"
fi

api="apigeeregistry.googleapis.com"

echo STEP 1: Acquiring gcloud auth token...
gcloud config set project "${project_id}"
token="$(gcloud auth print-access-token)"

echo STEP 2: Getting Project number...
project_num="$(gcloud projects describe "${project_id}" --format="value(projectNumber)")"

echo STEP 3: Enabling Apigee Registry API...
gcloud services enable apigeeregistry.googleapis.com --project="${project_id}"

echo STEP 4: Enabling Key Management Service API...
gcloud services enable cloudkms.googleapis.com --project="${project_id}"

echo STEP 5: Enabling Service Usage API...
gcloud services enable serviceusage.googleapis.com --project="${project_id}"

echo STEP 6: Granting roles/apigeeregistry.admin permission to the user...
gcloud projects add-iam-policy-binding "${project_id}" \
--member user:"${user}" \
--role roles/apigeeregistry.admin

echo STEP 7: Creating Apigee Registry P4SA...
gcloud beta services identity create --service="${api}" --project="${project_id}"

p4sa=service-${project_num}@gcp-sa-apigeeregistry.iam.gserviceaccount.com
cmek_key_name="projects/${project_id}/locations/${location}/keyRings/${key_ring_name}/cryptoKeys/${key_name}"

echo STEP 8: Creating encryption key "$cmek_key_name"...
gcloud kms keyrings create "${key_ring_name}" \
--location "${location}" \
--project "${project_id}"
gcloud kms keys create "${key_name}" \
--keyring "${key_ring_name}" \
--location "${location}" \
--purpose "encryption" \
--project "${project_id}"

echo STEP 9: Granting Apigee Registry P4SA permission on encryption key...
gcloud kms keys add-iam-policy-binding "${key_name}" \
--location "${location}" \
--keyring "${key_ring_name}" \
--member serviceAccount:"${p4sa}" \
--role roles/cloudkms.cryptoKeyEncrypterDecrypter \
--project "${project_id}"

echo STEP 10: Triggering the long running provisioning API. See response below...
curl --request POST \
https://"${api}"/v1/projects/"${project_id}"/locations/"${location}"/instances?instance_id=default \
--header "Authorization: Bearer ${token}" \
--header 'Content-Type: application/json' \
--data-raw '{
  "config": {
    "cmek_key_name": "'"${cmek_key_name}"'"
  }
}'

echo The operation takes about 30-40 mins to finish. Check the response from STEP 10 to get the operation ID.