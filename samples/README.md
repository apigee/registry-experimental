# Integration examples with Apigee Registry

Here are some examples of how to integrate with Apigee Registry.

> The examples are meant to demonstrate ways to integrate 
> **Apigee Registry** with your existing tools

* **Registry Spec Renderer** 
  * It provides a renderer for the specs in your registry.
  * Renderers supported : Swagger UI, GraphiQL, AsyncApi Renderer. 

* **Registy Mock Generator** 
  * It provides the ability to generate mock endpoints using prism.
  * This sample demonstrates the use of the Controller framework provided 
    by the Apigee Registry project.
  * A custom controller checks the state of the Registry and deploys the prism 
    mock service for every OpenAPI specification. 
  * An artifact (`**_mock-prism-endpoint_**`) is generated, once the mock 
    service is deployed successfully.
