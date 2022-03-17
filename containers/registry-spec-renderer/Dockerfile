FROM launcher.gcr.io/google/nodejs
ENV PORT=80
ENV APG_REGISTRY_ADDRESS=apigeeregistry.googleapis.com:443

COPY . /app/

# Install dependencies.
RUN npm --unsafe-perm install