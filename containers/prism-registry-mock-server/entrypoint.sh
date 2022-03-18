#!/usr/bin/env sh

registry get $REGISTRY_SPEC --contents > /openapi.yaml

prism mock -h "0.0.0.0" /openapi.yaml