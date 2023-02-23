# registry-connect

`registry-connect` is a command-line tool for extracting proxy data from an
Apigee X instance to be applied onto an API Registry.

## Usage

Running one of these commands retrieves either products or proxies, respectively,
from an Apigee X or SaaS runtime instance and formats it as API Registry-compatible
YAML:

    registry-connect discover apigee products ORGANIZATION
    registry-connect discover apigee proxies ORGANIZATION

The output from this command can be piped to `registry apply -` like so:

    registry-connect discover apigee products ORGANIZATION | registry apply -

Alternatively, the output may be sent to a file for inspection or processing,
at which point `registry apply -f FILE` can be run against it to apply it
to the registry. Example:

    registry-connect discover apigee products ORGANIZATION > apigee-apis.yaml
    registry apply -f apigee-apis.yaml

See `registry apply --help` for more information.

## Authentication

`registry-connect` uses
[Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
to connect to Apigee X. These are stored in your local environment when you login with `gcloud`:

`gcloud auth application-default login`

**MacOS note:** To run the `registry-connect` tool on MacOS, you may need to
[unquarantine](https://discussions.apple.com/thread/3145071) it by running the
following on the command line:

    xattr -d com.apple.quarantine registry-connect
