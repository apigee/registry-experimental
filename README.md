[![Go Actions Status](https://github.com/apigee/registry-experimental/workflows/Go/badge.svg)](https://github.com/apigee/registry-experimental/actions)

# Registry Experiments

The [Apigee Registry API](https://github.com/apigee/registry) allows teams to
upload and share machine-readable descriptions of APIs that are in use and
development.

This repository holds experimental code that builds on the Registry API: new
projects, suspended projects, and work that might be useful in future projects.

### registry-experimental

[cmd/registry-experimental](cmd/registry-experimental) contains a command-line
tool that is structurally identical to the registry tool but containing
experimental capabilities, some of which might eventually be migrated to the
registry tool.

### Running the registry-graphql proxy

[cmd/registry-graphql](cmd/registry-graphql) contains a simple proxy that
provides a read-only GraphQL interface to the Registry API. It can be run with
a local or remote `registry-server`.

### Running the registry-server and viewer using docker

We publish the Registry Server, Registry Viewer and Registry Envoy proxy docker
images to GitHub Packages. We have published a sample
[docker-compose.yml](docker-compose.yml) file to make is easier to set this up
on your local environment.

Steps to run this setup:

1. Create a Google oAuth Client ID.

   Use the following values for the
   [form](https://console.cloud.google.com/apis/credentials/oauthclient):

   - **Select the correct GCP Project.**
   - Application Type : Web Application
   - Name : API Registry
   - Javascript Authorized origins : http://localhost:8888
   - Authorized redirect URIs: http://localhost:8888

2. Update the `GOOGLE_SIGNIN_CLIENTID` value, in `docker-compose.yml` file,
   with the client ID from previous step.

3. Run the `docker compose up` command to start the server

4. Setup environment variables to use registry tools
   ```shell
       export APG_REGISTRY_ADDRESS="localhost:9999"
       export APG_ADMIN_ADDRESS=$APG_REGISTRY_ADDRESS
       export APG_ADMIN_INSECURE=1
       export APG_REGISTRY_INSECURE=1
   ```
5. Sample commands to create a project and api.
   ```shell
       apg admin create-project --project_id=project1
       apg registry create-api --api_id=api1 --parent=projects/project1/locations/global
       apg registry list-apis --parent=projects/project1/locations/global --json
   ```
6. To wipe out the local setup run `docker compose down -v`

## License

This software is licensed under the Apache License, Version 2.0. See
[LICENSE](LICENSE) for the full license text.

## Disclaimer

This is not an official Google product. Issues filed on Github are not subject
to service level agreements (SLAs) and responses should be assumed to be on an
ad-hoc volunteer basis.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING](CONTRIBUTING.md) for notes
on how to contribute to this project.
