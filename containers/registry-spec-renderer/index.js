/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
const express = require('express')
const fs = require("fs");
const handlebars = require("handlebars");
const {RegistryClient} = require("@google-cloud/apigee-registry");
const {Storage} = require('@google-cloud/storage');
const {credentials} = require("@grpc/grpc-js");
const jsYaml = require('js-yaml');
const {parseURL} = require("whatwg-url");
const cors = require('cors')
const storage = new Storage();

var client_options = {};
if (process.env.APG_REGISTRY_INSECURE
    && process.env.APG_REGISTRY_INSECURE == "1") {
  client_options.sslCreds = credentials.createInsecure();
}

const OPENAPI_MOCK_ENDPOINT = process.env.OPENAPI_MOCK_ENDPOINT;
const GRAPHQL_MOCK_ENDPOINT = process.env.GRAPHQL_MOCK_ENDPOINT;
const GRPC_DOC_ARTIFACT_NAME = process.env.GRPC_DOC_ARTIFACT_NAME
    || 'grpc-doc-url';

if (process.env.APG_REGISTRY_ADDRESS) {
  items = process.env.APG_REGISTRY_ADDRESS.split(":");
  client_options.apiEndpoint = items[0];
  client_options.port = items.length >= 1 ? items[1] : 443;
}

const client = new RegistryClient(client_options);

// Read the templates files for various api formats.
const swagger_ui_template = fs.readFileSync(
    "renderers/swagger-ui.html.template");
const graphiql_template = fs.readFileSync("renderers/graphiql.html.template");
const async_template = fs.readFileSync("renderers/async-ui.html.template");

const app = express();
app.set('trust proxy', true);
app.use(cors())
app.use(express.json());

//Serve static assets for openapi renderer
app.use('/swagger-ui',
    express.static(require('swagger-ui-dist').absolutePath()))

//Serve static assets for asyncapi renderer
app.use('/asyncapi/webcomponents/webcomponentsjs',
    express.static('node_modules/@webcomponents/webcomponentsjs'))
app.use('/asyncapi/web-component',
    express.static('node_modules/@asyncapi/web-component'))
app.use('/asyncapi/react-component',
    express.static('node_modules/@asyncapi/react-component'))

//Serve static assets for graphql renderer
app.use('/graphql/react', express.static('node_modules/react'))
app.use('/graphql/react-dom', express.static('node_modules/react-dom'))
app.use('/graphql/graphiql', express.static('node_modules/graphiql'))

function renderTemplate(res, apiFormat, spec_name) {
  var renderer_template = "";
  var api_endpoint = res.locals.endpoint_uri ? res.locals.endpoint_uri : "";
  switch (apiFormat) {
    case "openapi":
      renderer_template = swagger_ui_template.toString();
      if (!api_endpoint && OPENAPI_MOCK_ENDPOINT) {
        api_endpoint = OPENAPI_MOCK_ENDPOINT + "/" + spec_name;
      }
      break;
    case "asyncapi":
      renderer_template = async_template.toString();
      break;
    case "graphql":
      renderer_template = graphiql_template.toString();
      if (!api_endpoint && GRAPHQL_MOCK_ENDPOINT) {
        api_endpoint = GRAPHQL_MOCK_ENDPOINT + "/" + spec_name
      }
      break;
    case "grpc":
      renderer_template = "grpc_markdown";
      break
    default:
      renderer_template = "spec_file"
      break;
  }
  spec_url = "/spec/" + apiFormat + "/" + spec_name + (api_endpoint
      ? "?endpoint_uri=" + encodeURI(api_endpoint) : "");
  if (renderer_template == 'grpc_markdown') {
    /**
     * For GRPC documentation, we expect the HTML markup file to be
     * generated and stored in GCS. The URL to the GCS object will be
     * stored as an artifact on the spec object.
     */
    client.getArtifact({
      name: spec_name + '/artifacts/' + GRPC_DOC_ARTIFACT_NAME
    }, (err, artifact) => {
      client.getArtifactContents({
        name: spec_name + '/artifacts/' + GRPC_DOC_ARTIFACT_NAME
      }, async (err, artifact_content) => {
        if (err) {
          res.sendStatus(500);
          res.send("Error retrieving documentation for " + spec_name);
          res.end();
        } else {
          let artifact_url = artifact_content.data.toString().trim();
          let parsedUrl = parseURL(artifact_url);
          if (parsedUrl.host == 'storage.googleapis.com') {
            let bucket = parsedUrl.path.shift();
            let contents = await storage.bucket(bucket).file(
                parsedUrl.path.join("/")).download();
            res.setHeader("content-type", "text/html");
            res.send(contents[0].toString()).end();
          } else {
            res.redirect(artifact_url);
          }
        }
      })
    });

  } else if (renderer_template != "spec_file") {
    res.setHeader("content-type", "text/html; charset=UTF-8");
    hbstemplate = handlebars.compile(renderer_template);
    res.send(hbstemplate({specUrl: spec_url, apiEndpoint: api_endpoint}));
    res.end();
  } else {
    res.redirect(spec_url);
    res.end();
  }
}

function getAPIFormat(text) {
  let apiFormat = '';
  if (text.includes("openapi")) {
    apiFormat = 'openapi';
  } else if (text.includes("asyncapi")) {
    apiFormat = 'asyncapi';
  } else if (text.includes("discovery")) {
    apiFormat = 'discovery';
  } else if (text.includes("protobuf")) {
    apiFormat = 'grpc';
  } else if (text.includes("graphql")) {
    apiFormat = 'graphql';
  }
  return apiFormat;
}

/**
 * Render the spec associated to the Deployment.
 * From the deployment find out the spec revision to render.
 */
app.get(
    '/render/projects/:projectId/locations/:locationId/apis/:apiId/deployments/:deploymentId',
    (req, res) => {
      deployment_name = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId
          + "/deployments/" + req.params.deploymentId;
      client.getApiDeployment({
        name: deployment_name
      }, (err, deployment) => {
        if (err || !deployment.apiSpecRevision) {
          if (err) {
            console.error(err);
          } else {
            console.error(deployment_name + "not found");
          }
          res.sendStatus(500);
          res.end();
        } else {
          let queryString = "";
          if (deployment.endpointUri) {
            queryString += "endpoint_uri=" + encodeURI(deployment.endpointUri);
          }
          res.redirect(302,
              "/render/" + deployment.apiSpecRevision + (queryString ? "?"
                  + queryString : ""));
        }
      })
    });

/**
 * Render an API Spec from registry.
 * Chooses the renderer to use based on mimeType of the Spec
 */
app.get(
    '/render/projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId',
    (req, res) => {
      api_name = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId;
      spec_name = api_name + "/versions/"
          + req.params.versionId + "/specs/" + req.params.specId;
      if (req.query.endpoint_uri) {
        res.locals.endpoint_uri = decodeURI(req.query.endpoint_uri);
      }
      client.getApiSpec({
        name: spec_name
      }, (err, response) => {
        if (err) {
          console.error(err);
          res.sendStatus(500);
          res.end();
        } else {

          apiFormat = getAPIFormat(response.mimeType);

          if (!apiFormat) {
            client.getApi({
              name: api_name
            }, (err, response2) => {
              if (err) {
                console.error(err);
                res.sendStatus(500);
                res.end();
              } else {

                apiFormat = '';
                if (response2.labels['apihub-style']) {
                  apiFormat = getAPIFormat(
                      response2.labels['apihub-style']);
                }
                renderTemplate(res, apiFormat, spec_name);
              }
            });
          } else {
            renderTemplate(res, apiFormat, spec_name);
          }
        }
      })
    });

/**
 * Return the contents of an API Spec
 */
app.all(
    "/spec/:specType/projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId",
    (req, res) => {
      if (req.method !== 'POST' && req.method !== 'GET') {
        res.sendStatus(404);
        res.end();
        return;
      }
      if (req.query.endpoint_uri) {
        res.locals.endpoint_uri = decodeURI(req.query.endpoint_uri);
      }
      let spec_url = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId + "/versions/"
          + req.params.versionId + "/specs/" + req.params.specId;
      let api_url = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId;

      client.getApiSpecContents(
          {
            name: spec_url
          }, async (err, response) => {
            if (err) {
              console.error(err);
              res.sendStatus(500);
              res.end();
            } else {
              specObj = {};
              if (req.params.specType == 'openapi') {
                try {
                  specObj = JSON.parse(response.data);
                } catch {
                  specObj = jsYaml.load(response.data);
                }
                return addMockAddressForOpenAPI(specObj, spec_url, api_url,
                    res);
              } else {
                res.setHeader("content-type", response.contentType);
                res.send(response.data);
                res.end();
              }
            }
          })
    });

/**
 * Filter the list of deployments with the matching spec revision.
 *
 * Add the endpointUri for the matched deployments to the servers object
 * for the API Specs
 *
 * @param openAPISpecObj
 * @param spec_url
 * @param res
 */
function addMockAddressForOpenAPI(openAPISpecObj, spec_url, api_url, res) {
  client.listApiDeployments({
    parent: api_url,
    filter: "api_spec_revision == '" + spec_url + "'"
  }, (err, deployments) => {

    if (err) {
      /**
       * Ignore errors since we will specify a mock endpoint
       * even if we cannot get list of deployments.
       */
      console.error(err);
    }

    if (!deployments) {
      deployments = [];
    }

    if (OPENAPI_MOCK_ENDPOINT) {
      deployments.push({
        endpointUri: OPENAPI_MOCK_ENDPOINT + "/" + spec_url,
        displayName: "Mock Service"
      });
    }

    /**
     * Handle OpenAPI 2.0 specs
     */
    if (openAPISpecObj.swagger && openAPISpecObj.swagger.startsWith("2")) {
      let endpoint = deployments.length > 0 ? deployments[0].endpointUri : "";
      if (res.locals.endpoint_uri) {
        endpoint = res.locals.endpoint_uri;
      }
      if (endpoint) {
        let parsedUrl = parseURL(endpoint);
        openAPISpecObj.host = parsedUrl.host + (parsedUrl.port ? ":"
            + parsedUrl.port : "");
        openAPISpecObj.basePath = "/" + parsedUrl.path.join("/");
        openAPISpecObj.schemes = [parsedUrl.scheme];
      }
    }
    deployments.forEach(deployment => {
      if (!deployment.endpointUri) {
        return;
      }
      if (openAPISpecObj.openapi &&
          openAPISpecObj.openapi.startsWith("3")) {
        openAPISpecObj = openAPISpecObj || {
          servers: []
        }
        openAPISpecObj.servers.push(
            {
              url: deployment.endpointUri,
              description: deployment.displayName
            });
      }
    })

    res.json(openAPISpecObj);
    res.end();
  });
}

app.get('/', (req, res) => {
  res.sendStatus(200);
  res.end();
});

app.listen(process.env.PORT || 80, function (err) {
  if (err) {
    console.log("Error in server setup")
  }
  console.log("Server listening on Port", process.env.PORT || 80);
});
