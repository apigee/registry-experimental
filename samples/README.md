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

## Uploading to the hosted instance:
* Upload the whole directory using the following command:

```
registry apply -f googleapis -R --parent projects/$PROJECT/locations/global
```

where `PROJECT` is set to the name of the target project.

Note: Please make sure you are running a registry instance of version `v0.5.3` or higher.

* Wait for 5 to 15 minutes to see the artifacts to show up in the registry. 

## Uploading to a locally running instance:
* Upload the whole directory using the following command:

```
registry apply -f googleapis -R --parent projects/$PROJECT/locations/global
```

where `PROJECT` is set to the name of the target project.

Note: Please make sure you are running a registry instance of version `v0.5.3` or higher.

* Now, since you are running the local version, you will need to manually invoke the controller to generate the expected artifacts.
*  Run the following command to invoke the controller:
```
registry resolve projects/$PROJECT/locations/global/artifacts/apihub-manifest
```

Note: Please make sure you have installed [spectral linter](https://meta.stoplight.io/docs/spectral/ZG9jOjYyMDc0Mw-installation) version 5.X.X

* Due to the dependency relationships between different artifacts, you will need to invoke the controller at least 3 times in order to get all the artifacts generated. The artifacts generated in each invocation will be as follows:
    - First invocation:
        - vocabulary (you might see some errors here due to invalid spec definitions)
        - complexity (you might see some errors here due to invalid spec definitions)
        - conformance-apihub-styleguide (stores conformance reports)
        - connformance-receipt (stores the receipt of conformance action)
    - Second invocation:
        - score-lint-warnings (stores the score representing number of lint warnings)
        - score-lint-errors (stores the score representing number of lint errors)
        - score-receipt (stores the receipt of score computation action)
    - Third invocation:
        - scorecard-lint-summary (stores the  summary of all lint scores)
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

registry apply -f googleapis -R --parent projects/shrutiparab-sandbox/locations/global