# apigee-export

`apigee-export` is a command-line tool for exporting information from Apigee
into the API Registry. Apigee resources are exported as YAML files that can be
imported into API Registry using the `registry` tool.


**MacOS note:** To run the `apigee-export` tool on MacOS, you may need to
[unquarantine](https://discussions.apple.com/thread/3145071) it by running the
following on the command line:

```
xattr -d com.apple.quarantine apigee-export
```
