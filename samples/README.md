# Sample projects

Each subdirectory contains a sample project that can be uploaded using the
`registry` tool, e.g.

## googleapis
This directory contains yamls which will upload ~60 APIs to the registry, aalong with some project level configurations which will calculate the following artifacts for your specs:
- vocabulary
- complexity
- connformance reports
- scores
- scorecards 

* Upload the whole directory using the following command:

```
registry apply -f googleapis -R --parent projects/$PROJECT/locations/global
```

where `PROJECT` is set to the name of the target project.

Note: Please make sure you are running a registry instance of version `v0.5.3` or higher.

* Wait for 5 to 15 minutes to see the artifacts to show up in the registry. You can use the following commands to list the artifacts
```
# Listing vocabulary artifacts
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/vocabulary

# Listing complexity artifacts
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/complexity

# Listing conformance artifacts
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/conformance-apihub-styleguide

# Listing score artifacts
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/score-lint-errors
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/score-lint-warnings

# Listing scorecard artifacts
registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/scorecard-lint-summary
```
