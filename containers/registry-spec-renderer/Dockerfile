FROM launcher.gcr.io/google/nodejs

ENV PORT=8080
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443

COPY . /app/

RUN npm --unsafe-perm install && rm /app/.npmrc