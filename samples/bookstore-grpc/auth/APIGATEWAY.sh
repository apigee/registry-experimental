#!/bin/sh

PROJECT=$(gcloud config get project)
ADDRESS=$(gcloud api-gateway gateways describe bookstore --location=us-west2 --format "value(defaultHostname)")

export BOOKSTORE_BOOKSTORE_ADDRESS=$ADDRESS:443
unset BOOKSTORE_BOOKSTORE_INSECURE

export BOOKSTORE_BOOKSTORE_TOKEN=invalid
