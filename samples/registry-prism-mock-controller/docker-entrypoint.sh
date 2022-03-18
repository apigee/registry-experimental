#!/bin/sh

set -ex

registry upload manifest /prism-manifest.yaml --project-id=${REGISTRY_PROJECT_NAME} || true

registry resolve projects/${REGISTRY_PROJECT_NAME}/locations/global/artifacts/apihub-prism-mocker-manifest

rc=$(echo $?)

curl -fsI -X POST http://localhost:15020/quitquitquit

exit $rc
