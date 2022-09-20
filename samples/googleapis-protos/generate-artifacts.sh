#!/bin/bash

PROJECT=googleapis

guidelines=$(registry get projects/$PROJECT/locations/global/artifacts/apihub-styleguide --contents | jq -r .guidelines[].id)

scorePatterns=""

for g in $guidelines
do
    gid=${g:3}
    echo "apiVersion: apigeeregistry/v1
kind: ScoreDefinition
metadata:
  name: apihub-lint-$gid
data:
  displayName: Lint Errors
  description: Number of lint errors found for AIP-$gid in conformance report
  uri: https://linter.aip.dev/$gid
  uriDisplayName: AIP-$gid rules
  targetResource:
    pattern: apis/-/versions/-/specs/-
    filter: \"\"
  scoreFormula:
    artifact:
      pattern: \$resource.spec/artifacts/conformance-apihub-styleguide
      filter: \"\"
    scoreExpression: 'has(guidelineReportGroups[2].guidelineReports) ? sum(guidelineReportGroups[2].guidelineReports.filter(g, g.guidelineId == \"$g\").map(r, has(r.ruleReportGroups[1].ruleReports) ? size(r.ruleReportGroups[1].ruleReports) : 0)) : 0'
    referenceId: \"\"
  integer:
    minValue: 0
    thresholds: []" > apihub-lint-$gid.yaml
    scorePatterns+="\n    - \$resource.spec/artifacts/score-apihub-lint-$gid" 
done

echo -e "apiVersion: apigeeregistry/v1
kind: ScoreCardDefinition
metadata:
  name: apihub-lint-aip-summary
data:
  displayName: Lint Summary
  description: Summary of aip lint scores
  targetResource:
    pattern: apis/-/versions/-/specs/-
    filter: \"\"
  scorePatterns:$scorePatterns" > apihub-lint-aip-summary.yaml

