#!/bin/sh

PROJECT=$(gcloud config get project)
ADDRESS=$(gcloud run services describe bookstore --format "value(status.url)")
ADDRESS=${ADDRESS#https://}

export BOOKSTORE_BOOKSTORE_ADDRESS=$ADDRESS:443
unset BOOKSTORE_BOOKSTORE_INSECURE

export BOOKSTORE_BOOKSTORE_TOKEN=invalid
