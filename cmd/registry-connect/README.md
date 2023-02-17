# apigee-export

`apigee-export` is a command-line tool for extracting proxy data from an
Apigee X instance to be applied onto an API Registry.

## Usage

`apigee-export apis ORGANIZATION [DIRECTORY]`
`apigee-export deployments ORGANIZATION [DIRECTORY]`

Running this command exports Registry-compatible YAML files from an Apigee X
instance into the specified DIRECTORY. (If no DIRECTORY is specified, it will
only print to the console.)

Once `apigee-export` has been successfully run, the entire exported directory or
individual files can be imported into an API Registry instance by running:

`registry apply -f DIRECTORY|FILE`

See `registry apply --help` for more information.

## Authentication

`apigee-export` uses
[Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
to connect to Apigee X. These are stored in your local environment when you login with `gcloud`:

`gcloud auth application-default login`

**MacOS note:** To run the `apigee-export` tool on MacOS, you may need to
[unquarantine](https://discussions.apple.com/thread/3145071) it by running the
following on the command line:

```sh
xattr -d com.apple.quarantine apigee-export
```
