const express = require('express')
const fs = require("fs");
const handlebars = require("handlebars");
const {RegistryClient} = require("@giteshk-org/apigeeregistry");
const grpc = require("@grpc/grpc-js");

var client_options = {};
if (process.env.APG_REGISTRY_INSECURE
    && process.env.APG_REGISTRY_INSECURE == "1") {
  client_options.sslCreds = grpc.credentials.createInsecure();
}

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
    case "async":
      renderer_template = async_template.toString();
      break;
    case "graphql":
      renderer_template = graphiql_template.toString();
      break;
  }
  specUrl = "/spec/" + spec_name;
  if (renderer_template != "") {
    res.setHeader("content-type", "text/html; charset=UTF-8");
    hbstemplate = handlebars.compile(renderer_template);
    res.send(hbstemplate({specUrl: specUrl}));
  } else {
    res.redirect(specUrl);
  }
  res.end();
}

function discoverAPIFormat(text) {
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
          res.sendStatus(500);
          res.end();
        } else {
          apiFormat = discoverAPIFormat(response.mimeType);
          if (apiFormat == '') {
            client.getApi({
              name: api_name
            }, (err, response2) => {
              if (err) {
                res.sendStatus(500);
                res.end();
              } else {
                apiFormat = '';
                if (response2.labels['apihub-style']) {
                  apiFormat = discoverAPIFormat(
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
    "/spec/projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId",
    (req, res) => {
      if (req.method !== 'POST' && req.method !== 'GET') {
        res.sendStatus(404);
        res.end();
        return;
      }

      let spec_url = "projects/" + req.params.projectId + "/locations/"
          + req.params.locationId + "/apis/" + req.params.apiId + "/versions/"
          + req.params.versionId + "/specs/" + req.params.specId;

      client.getApiSpecContents(
          {
            name: spec_url
          }, (err, response) => {
            if (err) {
              console.error(err);
              res.sendStatus(500);
              res.end();
            } else {
              res.setHeader("content-type", response.contentType);
              res.send(response.data);
              res.end();
            }
          })
    });

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