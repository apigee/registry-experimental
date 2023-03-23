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
[docker-compose.yml](docker-compose.yml) file to make it easier to set this up
on your local environment.

Steps to run this setup:

1. Run the `docker compose up` command to start the server.

2. Point your [registry configuration](https://github.com/apigee/registry/wiki/registry-config)
   at your local instance.
   ```shell
       registry config configurations create local
       registry config set registry.address localhost:9999
       registry config set registry.insecure true
   ```

3. Sample commands to create a project and a sample API.
   ```shell
       registry rpc admin create-project --project_id=project1
       registry rpc create-api --api_id=api1 --parent=projects/project1/locations/global
       registry rpc list-apis --parent=projects/project1/locations/global --json
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
