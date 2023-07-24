# registry-bigquery

`registry-bigquery` is an experimental command-line tool that uses BigQuery to
index and search information in an API registry.

## Usage

Subcommands expect users to have a Google Cloud project with BigQuery enabled.
By default, this is assumed to be the same project as the configured registry
project, but this can be overridden with a command-line flag.

The following subcommands build indexes of information in API specs:

```
registry-bigquery index operations PATTERN
registry-bigquery index servers PATTERN
registry-bigquery index info PATTERN
```

## Authentication

`registry-bigquery` uses
[Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
to connect to BigQuery. These are stored in your local environment when you
login with `gcloud`:

`gcloud auth application-default login`

**MacOS note:** To run the `registry-bigquery` tool on MacOS, you may need to
[unquarantine](https://discussions.apple.com/thread/3145071) it by running the
following on the command line:

    xattr -d com.apple.quarantine registry-bigquery
