# Bleve

This directory contains an experimental search implementation built with
https://blevesearch.com. Specs are indexed as full-text blobs and queried with
the Bleve default queries.

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
