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
const {RegistryClient} = require("@giteshk-org/apigeeregistry");
const {credentials} = require("@grpc/grpc-js");
const jsYaml = require('js-yaml');
const { parseURL } = require("whatwg-url");
var cors = require('cors')

var client_options = {};
if (process.env.APG_REGISTRY_INSECURE
    && process.env.APG_REGISTRY_INSECURE == "1") {
  client_options.sslCreds = credentials.createInsecure();
}

const MOCK_ENDPOINT = process.env.MOCK_ENDPOINT;

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
app.use(cors())

//Serve static assets for openapi renderer
app.use('/renderer/openapi',
    express.static(require('swagger-ui-dist').absolutePath()))

//Serve static assets for asyncapi renderer
app.use('/renderer/async',
    express.static('node_modules/@webcomponents/webcomponentsjs'))
app.use('/renderer/async',
    express.static('node_modules/@asyncapi/web-component'))
app.use('/renderer/async',
    express.static('node_modules/@asyncapi/react-component'))

//Serve static assets for graphql renderer
app.use('/renderer/graphql', express.static('node_modules/react'))
app.use('/renderer/graphql', express.static('node_modules/react-dom'))
app.use('/renderer/graphql', express.static('node_modules/graphiql'))

function renderTemplate(res, apiFormat, spec_name) {
  var renderer_template = "";
  switch (apiFormat) {
    case "openapi":
      renderer_template = swagger_ui_template.toString();
      break;
    case "asyncapi":
      renderer_template = async_template.toString();
      break;
    case "graphql":
      renderer_template = graphiql_template.toString();
      break;
    default:
      renderer_template = "spec_file"
      break;
  }
  specUrl = "/spec/" + apiFormat + "/" + spec_name;
  if (renderer_template != "spec_file") {
    res.setHeader("content-type", "text/html; charset=UTF-8");
    hbstemplate = handlebars.compile(renderer_template);
    res.send(hbstemplate({specUrl: specUrl}));
  } else {
    res.redirect(specUrl);
  }
  res.end();
}

function getAPIFormat(text) {
  let apiFormat = '';
  if (text.includes("openapi")) {
    apiFormat = 'openapi';
  } else if (text.includes("asyncapi")) {
    apiFormat = 'asyncapi';
  } else if (text.includes("discovery")) {
    apiFormat = 'discovery';
  } else if (text.includes("proto")) {
    apiFormat = 'proto';
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
          res.redirect(302, "/render/" + deployment.apiSpecRevision);
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

      client.getApiSpec({
        name: spec_name
      }, (err, response) => {
        if (err) {
          console.error(err);
          res.sendStatus(500);
          res.end();
        } else {
          apiFormat = getAPIFormat(response.mimeType);
          if (apiFormat == '') {
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

      let spec_url = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId + "/versions/"
          + req.params.versionId + "/specs/" + req.params.specId;
      let api_url = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId;
      client.getApiSpecContents(
          {
            name: spec_url
          }, (err, response) => {
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

    if(!deployments) {
      deployments = [];
    }

    if (MOCK_ENDPOINT) {
      deployments.push({
        endpointUri: MOCK_ENDPOINT + "/" + spec_url,
        displayName: "Mock Service"
      });
    }

    deployments.forEach(deployment => {
      if (!deployment.endpointUri) {
        return;
      }
      if (openAPISpecObj.swagger && openAPISpecObj.swagger.startsWith("2")
          && !openAPISpecObj.host) {
        let parsedUrl = parseURL(deployment.endpointUri);
        openAPISpecObj.host = parsedUrl.host;
        openAPISpecObj.basePath = "/" + parsedUrl.path.join("/");
        openAPISpecObj.schemes = [parsedUrl.scheme];
      } else if (openAPISpecObj.openapi &&
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

app.get('/healthz', (req, res) => {
  res.sendStatus(200);
  res.end();
});

app.listen(process.env.PORT || 80, function (err) {
  if (err) {
    console.log("Error in server setup")
  }
  console.log("Server listening on Port", process.env.PORT || 80);
});
