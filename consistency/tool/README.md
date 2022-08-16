# API Consistency

API Vocabulary Consistency 

## Description


The Registry API allows teams to upload and share machine-readable descriptions of APIs that are in use and in development.

These API descriptions can be used by tools like linters, browsers, documentation generators, test runners, proxies, API client and server generators. The registry API has a supporting CLI tool that computes vocabulary used in API specs, the vocabulary tool. The consistency report tool builds on top of this; it analyzes vocabulary usage in specs, specifically how consistent usages are across APIs, versions, specs, etc. By consistent, we mean the uniformity of word definitions found in a project and specs associated with it.

Example: user_id, UserID, userid, USERID, etc are inconsistent usages of the same term.

3 main parts for the tool:
### Computing word groups:

Word groups are flavors of the same vocabulary. For example, user_id, UserID, userid, and USERID are word groups. There is a data struture to store inconsistent usages of vocabularies. The formal WordGroup definition and meaning terms can be found [here](https://github.com/apigee/registry-experimental/blob/main/google/cloud/apigeeregistry/applications/v1alpha1/consistency/word_group.proto).

```
message WordGroup {
   string id = 1 (REQUIRED);
   string kind = 2 (REQUIRED);
   map<string, int32> word_frequency = 3 (REQUIRED); 
  }
``` 
The longest phase of all, this step computes the word groups in an entire project in a periodic interval of time. 

### Computing Comparisons: 
Our reports are defined by the following data structure.

```
message ConsistencyReport {
    string id = 1 (REQUIRED);
    string  kind = 2 (REQUIRED);
    message Variation {
         string term = 1 (REQUIRED);
         WordGroup cluster = 2 (REQUIRED);
    }
    repeated Variation current_variations = 3;
    repeated string unique_terms = 5;
}
```
The ConsistencyReport formal definition and meaning of associated terms can be found [here](https://github.com/apigee/registry-experimental/blob/main/google/cloud/apigeeregistry/applications/v1alpha1/consistency/consistency_report.proto).

Given a spec name, this part of the tool creates a report about the word usage in the spec compared to the entire project. In this step, we detect word groups similar to each word found in the spec. We also detect unique words just introduced by a spec into the project.

### Interface: 
Once a comparison report has been generated for one (or several) specs, we generate an easily understood output csv file for users. In this step, the results we computed in the step above are finally provided in human-readable file(s). 

## Getting Started

Run 
```
pip install -r requirements.txt 
```
to install dependencies required by the tool from the root tool folder.

### Example usage

* First, set up a registry server running following the instructions [here](https://github.com/apigee/registry/blob/main/tests/demo/walkthrough.sh). Also make the protos defined in the `registry-experimental` folder from your virtual environment. To do so, run `make py-protos` from the root of `registry-experimental`.

* from the root of `registry`, run 
```
PROJECT=google
registry rpc admin create-project --project_id $PROJECT --json
```
* Save the `openapi-directory/APIs/googleapis.com ` folder somewhere.
* Upload the specs in the saved folder into the server using 
```
registry upload bulk openapi  --project-id $PROJECT {PATH TO FOLDER}         --base-uri https://github.com/APIs-guru/openapi-directory/blob/$COMMIT/APIs 

```
* Compute the vocabulary in the project using 
```
registry compute vocabulary projects/$PROJECT/locations/global/apis/-/versions/-/specs/-
```
* Now, to form WordGRoups, go to `tool/clustering` and run 
```
python3 main.py --project_name=google
```
Note, this will take a while as the project has more than 13000 words in it. If you are just interested in checking out how the tool runs e2e, change `line 53` of `clustering.py` from `self.words = valid_words` into something like `self.words = valid_words[1:1000]`. Adjust the number of words as needed. 

* Put the specs we have in the project into a path of you choice using 

```
registry list projects/google/locations/global/apis/-/versions/-/specs/- -> {YOUR PATH}/specs.txt
```

* From `tool/comparison` create a comparsion report for a spec of your choice from `specs.txt` using 

```
python3 main.py --project_name=google --spec_name={SPEC NAME OF YOUR CHOICE}
```
If you just want to just have a quick report use `python3 main.py --project_name=google --spec_name=projects/google/locations/global/apis/alertcenter/versions/v1beta1/specs/openapi.yaml`

* Generate a .csv file of the report using

```
python3 csv_generate.py --project_name=google --path={YOUR PATH HERE} --csv_name=google-report
```

* We can also use a shell script to query and give all the spec names we have to the report generator. Go to `tool/comparsion` again and run 
```
xargs -a {YOUR PATH TO spec.txt} -I{} -d'\n' python3 main.py --spec_name={} --project_name=google
```
This will generate a report for all the specs one by one. 

* Finally generate a report csv for all reports using
```
python3 csv_generate.py --project_name=google --path={YOUR PATH} --csv_name=google-report
```
