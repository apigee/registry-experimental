# Bleve

This directory contains an experimental search implementation built with
https://blevesearch.com. Specs are indexed as full-text blobs and queried with
the Bleve default queries.

Note that all calls below require that `registry-experimental` be configured to
use a `registry-server` instance.

Index specs with the following, where PATTERN should match one or more specs:

```
registry-experimental bleve index PATTERN
```

The index will be stored locally in `registry.bleve`. Use the `--bleve` option
to specify an alternate location.

Search the index with the following:

```
registry-experimental bleve search QUERY
```

Indexing and search are also available with a simple REST API that is provided
by `bleve serve`.

First run `registry-experimental bleve serve`. While it is running,
specs can be indexed and searched as follows:

Specs can be indexed by posting JSON to the `/index` endpoint:

```
curl http://localhost:8888/index \
    -X POST \
	-H "Content-Type: application/json" \
	-d @- \
	<<EOF
{
	pattern: "projects/${PROJECT_ID}/locations/global/apis/-/versions/-/specs/-",
	filter: "mime_type.contains('openapi')",
}
EOF

Note that the `filter` value is optional.

Specs be searched with `/search`:
```

curl http://localhost:8888/search?q=domain

```

This searches for specs containing the word "domain".
```
