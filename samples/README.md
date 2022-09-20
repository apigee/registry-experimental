# Sample projects

Each subdirectory contains a sample project that can be uploaded using the
`registry` tool, e.g.

## googleapis-openapi

This directory contains yamls which will upload ~60 APIs to the registry, along
with some project level configurations which will calculate the following
artifacts for your specs:

- connformance reports
- scores
- scorecards

## Uploading to the hosted instance:

- Upload the whole directory using the following command:

```
registry apply -f googleapis-openapi -R --parent projects/$PROJECT/locations/global
```

where `PROJECT` is set to the name of the target project.

Note: Please make sure you are running a registry instance of version `v0.5.3`
or higher.

- Wait for 5 to 15 minutes to see the artifacts to show up in the registry.

## Uploading to a locally running instance:

- Upload the whole directory using the following command:

```
registry apply -f googleapis-openapi -R --parent projects/$PROJECT/locations/global
```

where `PROJECT` is set to the name of the target project.

Note: Please make sure you are running a registry instance of version `v0.5.3`
or higher.

- Now, since you are running the local version, you will need to manually
  invoke the controller to generate the expected artifacts.
- Run the following command to invoke the controller:

```
registry resolve projects/$PROJECT/locations/global/artifacts/apihub-manifest
```

Note: Please make sure you have installed
[spectral linter](https://meta.stoplight.io/docs/spectral/ZG9jOjYyMDc0Mw-installation)
version 5.X.X

- Due to the dependency relationships between different artifacts, you will
  need to invoke the controller at least 3 times in order to get all the
  artifacts generated. The artifacts generated in each invocation will be as
  follows:
  - First invocation:
    - conformance-apihub-styleguide (stores conformance reports)
    - connformance-receipt (stores the receipt of conformance action)
  - Second invocation:
    - score-lint-warnings (stores the score representing number of lint
      warnings)
    - score-lint-errors (stores the score representing number of lint errors)
    - score-receipt (stores the receipt of score computation action)
  - Third invocation:
    - scorecard-lint-summary (stores the summary of all lint scores)
    - scorecard-receipt (stores the receipt of scorecard computation action)

## Verifying the results:

1. You can use the following commands to list the artifacts

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

2. To get the scores for a particular spec, run the following:

   ```
   # Get lint-errors
   registry projects/shrutiparab-sandbox/locations/global/apis/googleapis.com-analyticshub/versions/v1beta1/specs/openapi.yaml/artifacts/score-lint-errors --contents

   # Get lint-warnings
   registry get projects/shrutiparab-sandbox/locations/global/apis/googleapis.com-analyticshub/versions/v1beta1/specs/openapi.yaml/artifacts/score-lint-warnings --contents

   # Get scorecard
   registry get projects/shrutiparab-sandbox/locations/global/apis/googleapis.com-analyticshub/versions/v1beta1/specs/openapi.yaml/artifacts/scorecard-lint-summary --contents
   ```

3. The expected response should be something like follows:
   - score-lint-errors
     ```
     {
         "id": "score-lint-errors",
         "kind": "Score",
         "displayName": "Lint Errors",
         "description": "Number of lint errors found in conformance report",
         "definitionName": "projects/shrutiparab-sandbox/locations/global/artifacts/lint-errors",
         "integerValue": {
             "value": 1,
             "maxValue": 100
         }
     }
     ```
   - score-lint-warnings
     ```
     {
         "id": "score-lint-warnings",
         "kind": "Score",
         "displayName": "Lint Warnings",
         "description": "Number of lint warnings found in conformance report",
         "definitionName": "projects/shrutiparab-sandbox/locations/global/artifacts/lint-warnings",
         "integerValue": {
             "maxValue": 100
         }
     }
     ```
   - scorecard-lint-summary
     ```
     {
         "id": "scorecard-lint-summary",
         "kind": "ScoreCard",
         "displayName": "Lint Summary",
         "description": "Summary of lint scores",
         "definitionName": "projects/shrutiparab-sandbox/locations/global/artifacts/lint-summary",
         "scores": [
             {
             "id": "score-lint-errors",
             "kind": "Score",
             "displayName": "Lint Errors",
             "description": "Number of lint errors found in conformance report",
             "definitionName": "projects/shrutiparab-sandbox/locations/global/artifacts/lint-errors",
             "integerValue": {
                 "value": 1,
                 "maxValue": 100
             }
             },
             {
             "id": "score-lint-warnings",
             "kind": "Score",
             "displayName": "Lint Warnings",
             "description": "Number of lint warnings found in conformance report",
             "definitionName": "projects/shrutiparab-sandbox/locations/global/artifacts/lint-warnings",
             "integerValue": {
                 "maxValue": 100
             }
             }
         ]
     }
     ```

## How are Scores calculated?

- The calculation of Scores depends on the existence of ScoreDefintion
  artifacts. These artifacts are stored at the project level and hence are
  applied to the whole project.
- You can see what definitions currently exist in your registry using the
  following command:
  ```
  registry list projects/$PROJECT/locations/global/artifacts/- --filter='mime_type.contains("ScoreDefinition")'
  ```
- Every time the `compute score` command is invoked, it will fetch all the
  ScoreDefinitions in the project and compute Score artifacts based on the
  definitions which apply to the supplied spec.
- Here is an explanation of what each field means in a ScoreDefinition and how
  the `compute score` command uses it. Let's consider the ScoreDefinition that
  we are uploading as part of this sample.
  ```
  apiVersion: apigeeregistry/v1
  kind: ScoreDefinition
  metadata:
      name: apihub-lint-errors
  data:
      display_name: "Lint Errors"
      description: "Number of lint errors found in conformance report"
      uri: "https://meta.stoplight.io/docs/spectral/4dec24461f3af-open-api-rules"
      uri_display_name: "Spectral rules"
      target_resource:
          pattern: "apis/-/versions/-/specs/-"
      score_formula:
          artifact:
              pattern: "$resource.spec/artifacts/conformance-apihub-styleguide"
          score_expression: "has(guidelineReportGroups[2].guidelineReports) ? sum(guidelineReportGroups[2].guidelineReports.map(r, has(r.ruleReportGroups[1].ruleReports) ? size(r.ruleReportGroups[1].ruleReports) : 0)) : 0"
      integer:
          min_value: 0
          max_value: 100
  ```
  - **display_name**, **description**, **uri**, **uri_display_name**: These are
    metadata fields which can be set to attach additional information to the
    definition. In the hosted version, this information is used to populate the
    UI.
  - **target_resource**: this defines the resource pattern on which this score
    definition can be applied. Target resource is determined with the
    combination of `pattern` + `filter`.
    - **pattern**: This can be any valid resource pattern. In this definition
      we use "apis/-/versions/-/specs/-", which means this definition will be
      used to calculate `apihub-lint-errors` score for all the specs in the
      registry.
    - **filter**: This is a CEL filter which can be applied on the pattern. For
      example, if we had a filter which looked like
      `mime_type.contains("openapi")`, then our definition would have been used
      for spec of type `openapi` only.
  - **score_formula**: This defines how the score will be computed. It defines
    what the score depends on and also the formula to calculate the score.
    - **artifact**: artifact is a combination of `pattern` and `filter` which
      should fetch a single artifact, the contents from this artifact will be
      used to compute the score. In this example, we fetch the
      `conformance-apihub-styleguide` artifact which is attached to the target
      resource, note the use of `$resource.spec` to define that the artifact is
      under the target resource's spec. If the expected artifact doesn't exists
      or failure to fetch it will result in a failure in score calculation.
    - **score_expression**: This is a CEL expression which will be applied on
      the contents of the artifact to generate a score value. Any error in the
      expression will result in a failure in score calculation.
  - **rollup_formula**: This is another way of defining the formula which is
    essentially a wrapper around multiple `score_formulas`. See this
    [example](https://github.com/apigee/registry-experimental/blob/main/samples/googleapis-openapi/artifacts/apihub-lint-error-percentage.yaml#L12)
    for how it can be used.
  - **interger**: This defines the type of the generated score. The value
    generated by the `score_expression` should match the type defined here. In
    this case, the `score_expression` should always generate an integer value.
    - **min_value**: the lowest value this score can take
    - **max_value**: the highest value this score can take.
    - **thresholds**: you can also define thresholds which will be used to
      assign a Severity value to the score. Foe example on how to define
      thresholds, refer this
      [example](https://github.com/apigee/registry-experimental/blob/main/samples/googleapis-openapi/artifacts/apihub-lint-error-percentage.yaml#L24).
      The value that the `score_expression` generates should be within the min
      and max values defined here, if out of bounds, it is assigned a default
      Severity of `ALERT`
  - There is also an option to define
    [boolean](https://github.com/apigee/registry-experimental/blob/main/samples/googleapis-openapi/artifacts/apihub-lint-errors-exist.yaml#L16)
    or
    [percentage](https://github.com/apigee/registry-experimental/blob/main/samples/googleapis-openapi/artifacts/apihub-lint-error-percentage.yaml#L23)
    score types.
- For more details on the fields, refer the
  [proto](https://github.com/apigee/registry/blob/main/google/cloud/apigeeregistry/v1/scoring/definition.proto#L31)
  definition.

## How are ScoresCards calculated?

- The calculation of ScoreCards depends on the existence of ScoreCardDefintion
  artifacts. These artifacts are stored at the project level and hence are
  applied to the whole project.
- You can see what definitions currently exist in your registry using the
  following command:
  ```
  registry list projects/$PROJECT/locations/global/artifacts/- --filter='mime_type.contains("ScoreCardDefinition")'
  ```
- Every time the `compute scorecard` command is invoked, it will fetch all the
  ScoreDefinitions in the project and compute ScoreCard artifacts based on the
  definitions which apply to the supplied spec.
- Here is an explanation of what each field means in a ScoreCardDefinition and
  how the `compute scorecard` command uses it. Let's consider the
  ScoreCardDefinition that we are uploading as part of this sample.
  ```
  apiVersion: apigeeregistry/v1
  kind: ScoreCardDefinition
  metadata:
      name: lint-summary
  data:
      display_name: "Lint Summary"
      description: "Summary of lint scores"
      target_resource:
          pattern: apis/-/versions/-/specs/-
      score_patterns:
      - $resource.spec/artifacts/score-lint-errors
      - $resource.spec/artifacts/score-lint-warnings
      - $resource.spec/artifacts/score-lint-error-percentage
      - $resource.spec/artifacts/score-lint-errors-exist
  ```
  - **display_name**, **description**: These are metadata fields which can be
    set to attach additional information to the definition. In the hosted
    version, this information is used to populate the UI.
  - **target_resource**: this defines the resource pattern on which this
    ScoreCard definition can be applied. Target resource is determined with the
    combination of `pattern` + `filter`.
    - **pattern**: This can be any valid resource pattern. In this definition
      we use "apis/-/versions/-/specs/-", which means this definition will be
      used to calculate `apihub-lint-errors` score for all the specs in the
      registry.
    - **filter**: This is a CEL filter which can be applied on the pattern. For
      example, if we had a filter which looked like
      `mime_type.contains("openapi")`, then our definition would have been used
      for spec of type `openapi` only.
  - **score_patterns**: This definition is quite simple in the sense that you
    have to just list the score artifacts you want to be included in you
    ScoreCard. Note the use of `$resource.spec` to define that the artifact is
    under the target resource's spec. We are listing four artifacts here which
    will be wrapped in the lint-summary ScoreCard.
- For more details on the fields, refer the
  [proto](https://github.com/apigee/registry/blob/main/google/cloud/apigeeregistry/v1/scoring/definition.proto#L243)
  definition.
