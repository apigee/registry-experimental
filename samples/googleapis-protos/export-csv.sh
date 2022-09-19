#!/bin/sh

PROJECT=googleapis

artifacts=$(registry list projects/$PROJECT/locations/global/apis/-/versions/-/specs/-/artifacts/scorecard-apihub-lint-aip-summary)
regex="projects/[A-Za-z0-9-.]+/locations/global/(?<name>apis/[A-Za-z0-9-.]+/versions/[A-Za-z0-9-.]+/specs/[A-Za-z0-9-.]+)/artifacts/[A-Za-z0-9-.]+"
cols=""
rows=""
for a in $artifacts
do
    contents=$(registry get $a --contents)
    json=$(echo "{\"name\": \"$a\"} $contents" | jq -s add)

    if [ -z "$cols" ]
    then
        cols=$(echo $json | jq -r '["name", (.scores[].id | sub("score-apihub-lint-"; "")) ] | @csv')
    fi

    row=$(echo $json | jq -r --arg regex "$regex" '[(.name | capture($regex) | .name ), .scores[].integerValue.value] | @csv')
    rows+="$row\n"
done

# Output to a csv file
echo "$cols\n$rows" > aip-scores.csv
