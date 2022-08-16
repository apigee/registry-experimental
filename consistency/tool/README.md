# Project Title

API Vocabulary Consistency 

## Description


The Registry API allows teams to upload and share machine-readable descriptions of APIs that are in use and in development.

These API descriptions can be used by tools like linters, browsers, documentation generators, test runners, proxies, API client and server generators. The registry API has a supporting CLI tool that computes vocabulary used in API specs, the vocabulary tool. The consistency report tool builds on top of this; it analyzes vocabulary usage in specs, specifically how consistent usages are across APIs, versions, specs, etc. By consistent, we mean the uniformity of word definitions found in a project and specs associated with it.

Example: user_id, UserID, userid, USERID, etc are inconsistent usages of the same term.

3 main parts for the tool:
### Computing word groups:

Word groups are flavors of the same vocabulary. For example,user_id, UserID, userid, USERID are word groups. They are our data struture to store inconsistent usages of vocabularies. The WordGroup definition can be found here.

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
The ConsistencyReport definition can be found here.

Given a spec name, this part of the tool creates a “report” about the word usage in the spec compared to the entire project. In this step, we detect word groups similar to each word found in the spec. We also detect unique words just introduced by a spec into the project.

### Interface: 
Once a comparison report has been generated for one (or several) specs, we generate an easily understood output csv file for users. In this step, the results we computed in the step above are finally provided in human-readable file(s). 

## Getting Started

Run 
```
pip install -r requirements.txt 
```
to install dependencies required by the tool from the tool folder.

### Example

* How to run the program
* Step-by-step bullets
```
code blocks for commands
```

