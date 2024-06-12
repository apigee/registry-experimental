# Zero

Simple API management without a gateway. Instead of using a gateway, we will
make direct calls to service control APIs from within our API server.

Is this a good idea?

Pros:

- simple (no proxies to set up and manage)
- inexpensive (no additional sidecars to operate)
- runs anywhere (even when you run your server locally)

Cons:

- requires changes to your application
- hard to govern, there may be uncontrolled APIs that leak information

## Demonstration

### Preparation

#### You need a domain

To register a service with service manager, a domain is required.

##### Use a domain you control

We can register a domain with a registrar and prove to Google that we own it,
and then we can create services on that domain or any subdomain.
“example1.timbx.me”

##### Get a domain from App Engine

Alternately, we can use Google App Engine to get a domain that we can use. App
Engine apps are hosted at <appname>.appspot.com, where <appname> is usually the
project id.

- create an app engine app for your project
- this will give you “appname.appspot.com”. For example, for my project, named
  “nerdvana”, my domain name is “nerdvana.appspot.com”
- we can use this for our service name.
- we can also use subdomains of this domain.

#### Create OAuth credentials

These will be used to call the servicemanagement API

store them in ~/.config/zero/credentials.json

We could also use a service account for this.

#### Set up the CLI

create ~/.config/zero/zero.yaml for general configuration

```
serviceName: nerdvana.appspot.com
serviceConfig: 2023-10-06r1
apiKey: XXX-REDACTED-XXX
producerProject: nerdvana
consumerProject: nerdvana
summary: "Namaste"
title: "Nerdvana"
```

### Service Management

#### Create your service

Create your service with a call to the service management API

View the service in the endpoints console

#### Configure your service

Create your service config -- notice that we want to specify some things:

- name / version
- description
- operations

#### Rollout your configuration

Rollout your service config

Verify the rollout in the endpoints console

### Service Control

Create a service account to call the servicecontrol API

Call the check service

Call the check/allocatequota/report methods

### The Sample API

TODO

## Capabilities and Limitations

TODO

## Using the Service Management API to build an API Catalog

Get a list of managed services and the configurations for each. Each
configuration includes the list of operations that the API supports. These
lists can be used to build an index of APIs being provided by your Google
project.

Not all APIs -- these are just the formally registered ones. Other indicators
of APIs are Load Balancers, Cloud Run endpoints, and GKE ingresses.
