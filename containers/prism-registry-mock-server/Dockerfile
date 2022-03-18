FROM ghcr.io/apigee/registry-tools:latest as registry-tools
ENV REGISTRY_SPEC=""
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
ENV APG_REGISTRY_INSECURE=0

FROM launcher.gcr.io/google/nodejs
ENV REGISTRY_SPEC=""
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443
ENV APG_REGISTRY_INSECURE=0

EXPOSE 4010

COPY --from=registry-tools /bin/registry /bin/registry
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && npm install -g @stoplight/prism-cli

ENTRYPOINT ["/entrypoint.sh"]

