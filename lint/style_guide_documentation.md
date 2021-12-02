# API Style Guides in the Registry


## Introduction

One of the topics in the scope of governing APIs that is often overlooked is the importance of providing and enforcing concrete styling rules on API specs. There are many ways to write an API spec. Depending on the type of the spec (OpenAPI, AsyncAPI, Protobuf, etc.), it is easy to simply look up how to write a spec for your API, write it, and then forget about it. However, individuals and organizations often overlook the importance of enforcing good style in their API specs. The spec _specifies_ the API. It lays down the exact endpoints, descriptions, requirements, licenses, operations, and other metadata related to the API. Consequently, it is crucial to ensure the correctness of the API spec and to make sure that it is following industry-standard best practices.

The project I worked on during my internship at Google examines the importance of enforcing styling rules on API specs, and provides concrete data structures and a system through which organizations can declaratively enforce important styling rules on their API specs. The project in which my domain of work is used is the [open-source API Registry](https://github.com/apigee/registry) project on Github. We define an **API Style Guide** as the data structure through which organizations can declare style rules on their API specs. 


## Motivation


### Current Solution

To understand the motivation for an API style guide, we should understand the current solution for linting in the registry. Currently, the way linting is done is rudimentary. It consists of a single command compute lint, which lints any given resource based on its mime type with two major linters. If the mime type is OpenAPI, then [Spectral](https://stoplight.io/open-source/spectral/) is used. If the mime type is a proto, then [api-linter](https://github.com/googleapis/api-linter) is used. The types of rules that are used are just the common rules, and the user does not have much control over the types of rules that are used.


### AIPs (API Improvement Proposals)

The api-linter is a linter for AIPs. AIPs are defined on their website as “focused design documents for flexible API development.” They give users a lot of control over the design of API specs.

Each AIP contains a set of design rules that achieve a common design goal for API Protobuf specifications. [Here](https://google.aip.dev/126) is an example of an AIP that governs styling on enumeration types in Protobuf specs.


### Envisioning a solution for the API Registry

AIPs are intuitive, human readable, and a great form of documentation. They work very well for protos, then if we can turn them into a data structure that is machine readable, we could virtually use them for any type of API spec (OpenAPI, AsyncAPI, GraphQL, etc.). If we call this new format an API style guide, then we are designing a data structure that governs all API specs in the registry.

One major benefit of this new solution would be that it gives users the complete and utter control of their API specs’ design. They can enable rules that they want, deprecate them, change their status, and outline their severities. It is essentially a language to express the final state of the API specs, and the features that we will add will take care of the rest.


## User Journey

When designing a solution for the API Registry project, it was really important to keep the user journey in mind. How would the user interact with the system? What are the steps that the user could take to create, enforce, and gather results from API style guides? What if the user wanted to use their own linters rather than the ones we provide out-of-box? These were some of the questions that needed to be answered when brainstorming this stage. The following were some of the **functional requirements** from the user perspective:



* The user should be able to create a style guide to govern API specs in their organization
    * The style guide needs to be a_ declarative_, simple and intuitive data structure that is relatively easy to create (similar to how AIPs are a simple and intuitive documentation page). It should make sense to a user at first glance what the various fields mean, and how it can be created.
* The user should be able to upload the style guide to the registry
* The user should be able to get a report of how well different specs in their organization conform to the style guide.
* The user should be able to create and specify their own custom linters if they choose to do so, but should have access to the default rules provided by community-standard linters out-of-box.


## Non-Functional Requirements

The three main NFRs in this project that were kept in mind when designing this solution were **ease of use**, **reliability**, and **maintainability**.



* **Ease of use** is perhaps the most important non-functional requirement in this project. This feature is going to be used by our customers and by the open-source community. If we make the process of governing an API using style guides difficult, then people will likely opt for other approaches. We wanted to keep our interfaces simple, the commands simple, and creating custom linters simple. This requirement motivated many of the architectural decisions and the organization of our data structures that we present.
* **Reliability** is important for this project because we need to ensure that every request from the user is fulfilled successfully in a correct manner. If the organization needs to gather a report on every spec they have which violates a set of style rules, and we can’t provide that to them accurately, then our solution is ultimately a failure. We needed to make sure that all code that was being written kept reliability in mind, which meant thinking of corner-cases and various unexpected things the user might do and reporting the correct errors.
* **Maintainability** is very important in all software projects, and especially in open-source projects. There will always be modifications, fixes, enhancements, and features built on top of the code that is written, and thus the code needs to be extensible. The code was made maintainable, and various architecture decisions that were made kept maintainability in mind. There are also unit tests written for the various plugins that were created, which ensures that other engineers cannot make breaking changes.


## Project Goals

The goals of this project are summarized as follows. Many of these goals were established after researching stats on various APIs, brainstorming, and testing different solutions over the course of the internship.



* Create the API Style Guide data structure while keeping in mind the functional requirements, user journey, and NFRs.
* Create a conformance report data structure to report how well a spec conforms to an API Style Guide.
* Create a command that can compute the conformance reports for any given spec, according to all relevant style guides. \

* Create an architecture and framework whereby a user can create custom linters and specify them in their API Style Guide.
* Create the Spectral built-in plugin according to the architecture.
* Create the API-linter built-in plugin according to the architecture.
* Design and implement a proof-of-concept custom OpenAPI linter that lints important rules that Spectral does not provide out-of-box. This linter can be used as a reference by users to write their own custom linters.


## Data Structure: API Style Guide


### Introduction

An API Style Guide solves the problems that were mentioned earlier. They follow the conventions proposed by AIPs and turn them into a machine readable format. By allowing users to write and attach API Style Guides in the registry, we provide a powerful, declarative way to govern API specifications.


### Proto Definition 

To design the API Style Guide proto, we had to brainstorm an interface that would be intuitive for users to declare the rules that govern their API specs. As mentioned earlier, we also used [AIPs](http://aip.dev/) (API Improvement Proposals) as motivation when designing this solution.

The end result of what we created was as follows. An API Style Guide contains:



* **<span style="text-decoration:underline;">Guidelines</span>** that the API Spec should follow. These guidelines are general, and can be composed of multiple rules. Each guideline can be in one of four states: `PROPOSED, ACTIVE, DEPRECATED, DISABLED`. \

* Under each guideline, there is a set of **<span style="text-decoration:underline;">Rules</span>**. These rules are linter-specific, and so there is a field to specify which linter the rule belongs to. The Rule has a name that the user supplies, and the user is also required to provide the name of the rule on the linter which enforces it. Each rule can have one of four severities: `ERROR, WARNING, INFO, HINT`. \

* **<span style="text-decoration:underline;">Linter</span>** messages provide details regarding specific linters. API Style Guides are also used as a form of documentation. They should be able to be looked at and the governance of style on APIs should be understood. Linter messages summarize the name of linters used in the Style Guide, as well as a _uri_ field to link source code or documentation for the linter.

Here is the [permalink](https://github.com/apigee/registry/blob/69c1e47455dabbaa9e1b3aa0d10a90d004e29562/google/cloud/apigee/registry/applications/v1alpha1/registry_styleguide.proto) to the proto definition for API Style Guides in the registry.


### Scope of Governance

An API Style Guide can be attached to a project resource in the registry, and it should govern the style for all specs of a given mime type underneath that project. There can be multiple API Style Guides on a given project. In such a case, each individual API Style Guide may govern static checking on a different mime type, or govern specs of the same mime type while serving a different overall purpose.


### Linters and Rules

We provide out-of-the-box support for standard Spectral and API-Linter rules. Currently, Spectral lints AsyncAPI and OpenAPI specs, whereas API-Linter lints Protobuf. When specifying an API Style guide, you may specify `spectral` or `api-linter` under the `linter` field of a Rule, followed by a linter rule name.

For Spectral, linter rule names can be found at these links:



* OpenAPI: [https://meta.stoplight.io/docs/spectral/ZG9jOjExNw-open-api-rules](https://meta.stoplight.io/docs/spectral/ZG9jOjExNw-open-api-rules)
* AsyncAPI: [https://meta.stoplight.io/docs/spectral/ZG9jOjUzNDg-async-api-rules](https://meta.stoplight.io/docs/spectral/ZG9jOjUzNDg-async-api-rules)

For API Linter, they can be found here: [https://linter.aip.dev/rules/](https://linter.aip.dev/rules/)

We also provide the ability for users to create their own linters. This is [described in detail in a subsequent section](#bookmark=id.mfw2s0tgrdu) of this document.


### Creating and uploading a Style Guide to the Registry

API Style Guides can be specified in YAML format, and uploaded to the registry through the use of the upload subcommand, passing in the path to the style guide and the project ID to upload on.

Here is an example of an API Style Guide:


```yaml
name: openapistyleguide
mime_types:
  - application/x.openapi+gzip;version=2
guidelines:
  - name: refproperties
    display_name: Govern Ref Properties
    description: This guideline governs properties for ref fields on specs.
    rules:
      - name: norefsiblings
        description: An object exposing a $ref property cannot be further extended with additional properties.
        linter: spectral
        linter_rulename: no-$ref-siblings
        severity: ERROR
    status: ACTIVE
linters:
  - name: spectral
    uri: https://github.com/stoplightio/spectral
```


This style guide can be uploaded to the registry using this command:


```
registry upload styleguide <PATH TO STYLE GUIDE> --project_id <PROJECT ID>
```



## Data Structure: Conformance Report


### Introduction

Once a style guide has been created, we need to be able to get some sort of a report for each spec which specifies how _well_ the spec conforms to the style guide. To do this, we need the notion of a **conformance report **to each spec. The conformance report describes how well a spec conforms to a given style guide.


### Proto Definition 

To design a conformance report, we needed to determine what data format would be intuitive for an organization manager or general user to read and understand the Style Guide violations for any given spec. We decided that we first needed to separate the various Guideline violations according to their statuses. For instance, the user should be able to easily see only the Guidelines that were violated that had status `ACTIVE`, or only the status `PROPOSED`, etc. Furthermore, we decided that users should be able to separate various Rule violations _under_ each guideline based on severity. For example, the user should be able to see only rules that were violated with the status **ERROR**, or only rules with the status **WARNING**, etc.

What resulted was the following. A conformance report is composed of:



* A mapping between Guideline status to Guideline Reports. A Guideline Report summarizes all the rule violations of a specific Guideline. This allows us to fulfill the first requirement of being able to separate various Guideline violations based on status. \

* Each Guideline Report contains a mapping between severity to Rule Reports. A Rule  Report summarizes the violation of a specific rule on a linter. It gives the location where the rule was violated, a description of the rule, and a suggestion for fixing the violation. This also fulfills the second requirement of being able to separate Rule violations based on severity.

 

Here is the [permalink](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/google/cloud/apigeeregistry/applications/v1alpha1/registry_conformance_report.proto) for the definition of a conformance report in the registry.


## Computing Conformance Reports


### Introduction

So far, we have covered the following:



* What an API Style Guide is
* How to create an API Style Guide and upload it to the registry
* What a conformance report is

This section wraps all these concepts together and teaches how conformance reports can be computed for specs.


### Using the ‘compute conformance’ command to compute conformance reports

This functionality is provided by the conformance subcommand of registry compute. In specific, the following command:


```
registry compute conformance <SPEC RESOURCE NAME>
```


This command does the following tasks:

* Searches for all API Style Guides on the project that contains the specified spec
* For each API Style Guide “styleguide” on the project that supports the mime type of the spec:
    * Computes the conformance report for this spec for styleguide
    * Attaches the conformance report as an artifact on this spec, with the name of the artifact being “conformance-&lt;styleguide.name>”.

After the command executes, we provide the guarantee that all possible conformance reports have been computed for the spec that you ran the command on, provided there are no errors in the API Style Guides.


### Possible Extensions

One future extension for this command would be to add an extra flag that allows the user to specify a specific style guide or a range of style guides that the conformance report should be computed for. Currently, conformance reports are computed for all possible style guides that support the spec’s mime type, but this would provide granularity for the user.


### Source Code

The implementation details of the conformance subcommand will be described in a subsequent section, but the code can also be [found here](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/cmd/registry/cmd/compute/conformance.go).


## Facilitating Custom Linters


### Introduction

As mentioned earlier, the implementation of API Style Guides in the API Registry supports standard Spectral and API-Linter rules out of the box with no additional configuration required.

However, users may want to create their own linters, and we realized that this is something that we should facilitate in our implementation. As a result, we provide a simple way for users to create their own custom linters.


### Creating a custom linter

Creating a custom linter is simple: The user simply needs to create a binary executable that can receive a certain data format in standard-input (STDIN), and return a certain data format on standard-output (STDOUT).  From this point forward, we will call this binary executable a **plugin**. The plugin can be placed in the same directory as the registry which will allow it to be accessible. The name of the plugin should be `registry-lint&lt;LINTERNAME>”. This linter name will be the same one that will be specified in the rules of your style guide.

This follows an established pattern that both [gnostic](https://github.com/google/gnostic) and [protoc](https://github.com/protocolbuffers/protobuf) use – but this will be discussed later in implementation details.



* The input to the plugin will be called a `LinterRequest`. Its structure is described [here](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/google/cloud/apigeeregistry/applications/v1alpha1/registry_lint.proto#L120-L129).
* The output of the plugin will be called a `LinterResponse`. Its structure is described [here](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/google/cloud/apigeeregistry/applications/v1alpha1/registry_lint.proto#L135-L148).

To summarize, the input to the plugin is a directory that contains specs that need to be linted, along with a set of linter rules with which to lint. The linter responds with either a list of errors (empty if linting succeeded), or a `Lint` message containing Lint results.

An **example** that illustrates a sample plugin can be found [here](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/cmd/registry/plugins/registry-lint-sample/main.go). Also, here our implementations of the [Spectral](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/cmd/registry/plugins/registry-lint-spectral/main.go) and [API-Linter](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/cmd/registry/plugins/registry-lint-api-linter/main.go) plugin. The nice thing about the plugin architecture is that plugins can be implemented in any language. As long as your final executable takes in the input that we provide and returns the output that we expect, it will work.


## Implementation Details


### Introduction

This section will be of importance to Google engineers who are maintaining this code in the API Registry project, or anyone who is interested in extending any of the functionality that has been mentioned thus far. The implementation details will be discussed in sections that detail various parts of the overall project.


### Conformance Command

The conformance command takes in the resource name of a spec. This resource name can specify a range of specs as well via the wildcard (projects/-/apis/-/versions/-/specs/-). It determines the project that the provided specs belong to. The following is the detailed pseudocode for the conformance command:

Given a spec resource name (possibly specifying a range of specs via a wildcard):



* Find the project name of the spec resource
* For each style guide SG that is specified on the project:
    * If the mime type that SG governs contains the mime type of the spec resource:
        * Initialize a worker pool W of size 16.
        * For each spec S in the spec resource:
            * Spawn a thread in the worker pool that is responsible for computing the conformance of S with respect to SG
            * INSIDE THE THREAD:
            * Iterate through the SG to determine the linters that need to be used, and the rules that correspond to them.
            * For each linter L:
                    * Call the plugin corresponding to L with the rules that were determined.
            * Aggregate and organize all of the data returned from each linter into a conformance report
            * Attach the conformance report to S as an artifact


### Plugins

As mentioned earlier in this document, the registry makes use of a plugin-based design pattern to lint specs. How the plugins can be set up and created has already been documented earlier – this section will discuss the reasons for us choosing this design pattern, as well as the way the conformance command invokes these linters.



* **Language-agnostic**: Plugins can be written in any language provided they can be compiled into a binary, read the [provided format](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/google/cloud/apigeeregistry/applications/v1alpha1/registry_lint.proto#L120-L129) from STDIN and write the [expected format](https://github.com/apigee/registry/blob/bfaf5c2e3719e31f6de83225520d3d77035aa156/google/cloud/apigeeregistry/applications/v1alpha1/registry_lint.proto#L135-L148) to STDOUT. We don’t expect our users to understand Golang, and this isn’t something that they should be concerned with either.
* **Low Coupling, High Cohesion: **Before plugins, we were using another solution which required the user to implement a common linter interface. However, with this approach the main limitation was that the user would have to modify code in the Registry to register their new linter and write the code for their new linter. The user would also need to write Golang, and understand the Registry. With this new approach, the user does not need to understand the Registry at all, and does not need to modify any code. Their plugin can be written completely independently with all of its components internalized in the plugin (high cohesion), and does not need to interface with the registry at all (low coupling). As long as the user provides us with a binary executable with the right name, everything will work out.

This pattern isn’t something that we created from scratch without testing or without prior knowledge. It is an established pattern in both the protoc and gnostic projects by Google.

To invoke the linter plugin, the registry reads the name of the linter from the API Style Guide under Rule’s `linter` field. It then invokes the plugin with the name registry-lint-LINTERNAME. For example, suppose the name of the linter is spectral. The linter with the name registry-lint-spectral will be invoked. Details about how the data will be sent is covered in the [Conformance section of implementation details](#bookmark=id.mtr4ocnlx85).


### Extensions/Alternatives to Plugin Implementation

One viable extension that we considered was, alongside making linter plugins, to also provide a way for the user to create and register a gRPC service through which a linter can be executed. So rather than delegating the work of linting to a binary executable on disk, it would be sent to a web API. This is something that would be nice to have in the future.

One of the alternatives that was considered for the plugin implementation was to have a common interface for a linter that all linters implement. The user would simply implement this interface if they wanted to create a plugin, and register their linter in a factory method that would return the correct linter based on the name. We noticed that the main drawback of this approach was that it led to high coupling with the registry project. The user would actually have to modify the code. We decided not to opt for this in favor of the plugin implementation, but it is mentioned here for completeness.


### Extending Conformance Reports

Conformance reports are currently defined on each spec. However, it would be a good idea in the future to have conformance reports report a collection of specs. This would make for very interesting reports from the user perspective. For instance, if we attach a conformance report on an API, it could summarize the conformance for all the specs under that API. The conformance report proto would likely have to be modified to allow for some sort of aggregation so that reports are still a manageable size and readable, yet still summarize results for all specs underneath that resource.


## Testing


### Linter Plugins (Spectral, API Linter, Sample OpenAPI Linter)

All the linter plugins in this project were unit tested thoroughly.  The following are the permalinks to the unit tests for Spectral, API Linter, and the Sample OpenAPI Linter, respectively:



* [https://github.com/apigee/registry/blob/5795e43b8395c674ac24597d0d078c7a1514e61d/cmd/registry/plugins/registry-lint-spectral/main_test.go](https://github.com/apigee/registry/blob/5795e43b8395c674ac24597d0d078c7a1514e61d/cmd/registry/plugins/registry-lint-spectral/main_test.go)
* [https://github.com/apigee/registry/blob/main/cmd/registry/plugins/registry-lint-api-linter/main_test.go](https://github.com/apigee/registry/blob/main/cmd/registry/plugins/registry-lint-api-linter/main_test.go)
* [https://github.com/apigee/registry/blob/main/cmd/registry/plugins/registry-lint-openapi-sample/main_test.go](https://github.com/apigee/registry/blob/main/cmd/registry/plugins/registry-lint-openapi-sample/main_test.go)


### Upload Style Guide Command

The <code>upload</code> subcommand in the registry has an established testing pattern, which was followed for the <code>upload styleguide</code> command. The following is a permalink for the unit tests for the upload style guide command:</strong>

[https://github.com/apigee/registry/blob/5795e43b8395c674ac24597d0d078c7a1514e61d/cmd/registry/cmd/upload/styleguide_test.go](https://github.com/apigee/registry/blob/5795e43b8395c674ac24597d0d078c7a1514e61d/cmd/registry/cmd/upload/styleguide_test.go)**


### Compute Conformance Command

With regards to the conformance command: As it stands, the API Registry project does not have a uniform way to unit test or integration-test the command-line tools, and conventions surrounding unit testing are being cemented. A description of this can be found [here](https://github.com/apigee/registry/pull/296#issuecomment-920228806). As a result, I provide a manual method to test, until we create a method to test CLI tools as a whole in the registry.



* **Step 1**: Create a project and upload at least one spec into that project
* **Step 2**: Create an API Style Guide, and upload it to the project
* **Step 3**: Run the compute conformance command on a spec within the project that the style guide applies to (supports the mime type)
* **Step 4**: Get the contents of the conformance report that was generated

A concrete example of this is provided in the [following pull request](https://github.com/apigee/registry/pull/329).


## Appendix

**<span style="text-decoration:underline;">Example of a Protobuf Style Guide with API Linter</span>**


```yaml
name: google-aip
mime_types:
 - application/x.protobuf+zip
guidelines:
 - name: aip126
   display_name: Enumerations
   description: This guideline governs enum objects in proto files.
   rules:
     - name: upperSnakeCase
       description: All enum values must use UPPER_SNAKE_CASE.
       linter: api-linter
       linter_rulename: enumValueUpperSnakeCase
       severity: ERROR
     - name: unspecifiedSuffix
       description: >
         The first value of the enum should be the name of the enum itself
         followed by the suffix _UNSPECIFIED.
       linter: api-linter
       linter_rulename: unspecified
       severity: WARNING
 - name: aip140
   display_name: Field Names
   description: This guideline governs the style on field names in proto files.
   rules:
     - name: lowerSnakeCase
       description: Field definitions in protobuf files must use lower_snake_case names.
       linter: api-linter
       linter_rulename: lowerSnake
       severity: ERROR
     - name: uri_over_url
       description: Field names representing URLs or URIs should always use uri rather than url.
       linter: api-linter
       linter_rulename: uri
       severity: WARNING
   status: ACTIVE
linters:
 - name: api-linter
   uri: https://github.com/googleapis/api-linter
```

