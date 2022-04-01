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

import * as express from 'express';
import * as cors from 'cors';
import {RegistryClient} from '@giteshk-org/apigeeregistry';
import {ClientOptions} from 'google-gax';
import {credentials} from '@grpc/grpc-js';
const Yaml = require('js-yaml');

import {IHttpOperation} from '@stoplight/types';
import {Request, Response} from 'express';

const {
  getHttpOperationsFromSpec,
} = require('@stoplight/prism-cli/dist/operations');
const {
  createClientFromOperations,
} = require('@stoplight/prism-http/dist/client');

const app = express();
const PORT = process.env.PORT || 3000;
const HEADER_REGISTRY_SPEC =
  process.env.HEADER_REGISTRY_SPEC || 'apigee-registry-spec';

const client_options = <ClientOptions>{};

if (process.env.APG_REGISTRY_INSECURE) {
  if (process.env.APG_REGISTRY_INSECURE === '1') {
    client_options.sslCreds = credentials.createInsecure();
  }
}

if (process.env.APG_REGISTRY_ADDRESS) {
  const items = process.env.APG_REGISTRY_ADDRESS.split(':');
  client_options.apiEndpoint = items[0];
  client_options.port = items.length >= 1 ? parseInt(items[1]) : 443;
}

const client = new RegistryClient(client_options);

app.use(cors());

/**
 * Error Handling
 *
 * @param err
 * @param res
 */
function _sendError(err: Error, res: Response) {
  console.error('Error ' + err.message);
  res.status(500).send(new Error(err.message)).end();
}

/**
 * Process the mock request
 * This method will fetch the spec contents from Registry
 * and pass the spec to the Prism Library.
 *
 * Prism Library will generate responses based on the contents of the spec.
 * * Prism Documentation details of how responses are generated:
 * https://meta.stoplight.io/docs/prism/ZG9jOjk1-http-mocking#response-examples
 *
 * @param req
 * @param res
 */
function processMockRequest(req: Request, res: Response) {
  const spec_url =
    'projects/' +
    res.locals.projectId +
    '/locations/' +
    res.locals.locationId +
    '/apis/' +
    res.locals.apiId +
    '/versions/' +
    res.locals.versionId +
    '/specs/' +
    res.locals.specId;

  console.info(
    'Spec[' + spec_url + '] : ' + req.method + ' ' + res.locals.apiPath
  );
  //Return 200 to HEAD Requests
  if (req.method === 'HEAD') {
    res.setHeader('Access-Control-Allow-Origin', req.headers['origin'] || '*');
    res.setHeader(
      'Access-Control-Allow-Headers',
      req.headers['access-control-request-headers'] || '*'
    );
    res.setHeader('Access-Control-Allow-Credentials', 'true');
    res.setHeader(
      'Access-Control-Expose-Headers',
      req.headers['access-control-expose-headers'] || '*'
    );
    res.sendStatus(200);
    res.end();
    return;
  }

  client.getApiSpecContents(
    {
      name: spec_url,
    },
    (err, response) => {
      if (err) {
        return _sendError(err, res);
      } else {
        if (response && response.data) {
          const {data} = response;
          const specString = <string>data;
          let specObj: Object;
          try {
            specObj = JSON.parse(specString);
          } catch (e) {
            specObj = Yaml.load(specString);
          }
          // fs.writeFileSync(specObj, JSON.stringify(specObj));
          getHttpOperationsFromSpec(specObj)
            .then((operations: IHttpOperation[]) => {
              const prism = createClientFromOperations(operations, {
                // logger: createLogger('TestLogger'),
                mock: {dynamic: false},
              });
              prism
                .request(
                  res.locals.apiPath,
                  req
                  // operations
                )
                .then((output: any) => {
                  Object.keys(output.headers).forEach(key => {
                    res.setHeader(key, output.headers[key]);
                  });

                  res.status(output.status);
                  res.send(output.data);
                  res.end();
                })
                .catch((err: Error) => {
                  return _sendError(err, res);
                });
            })
            .catch((err: Error) => {
              return _sendError(err, res);
            });
        } else {
          return _sendError(new Error('Spec not found'), res);
        }
      }
    }
  );
}

function processParams(req: Request, res: Response) {
  for (const [key, value] of Object.entries(req.params)) {
    if (key === '0') {
      res.locals['apiPath'] = '/' + value;
    } else {
      res.locals[key] = value;
    }
  }
}

/**
 * Mock request endpoint with registry spec mentioned in the header
 */
app.all('/mock/*', (req: Request, res: Response) => {
  if (!req.headers[HEADER_REGISTRY_SPEC]) {
    return _sendError(
      new Error('Header "' + HEADER_REGISTRY_SPEC + '" not defined.'),
      res
    );
  }
  const parts = (<string>req.headers[HEADER_REGISTRY_SPEC]).split('/');
  if (parts.length !== 10) {
    return _sendError(
      new Error('Invalid value for ' + HEADER_REGISTRY_SPEC),
      res
    );
  }
  res.locals.projectId = parts[1];
  res.locals.locationId = parts[5];
  res.locals.apiId = parts[7];
  res.locals.specId = parts[9];
  processParams(req, res);
  processMockRequest(req, res);
});

/**
 * Mock request endpoint with registry spec in request Path
 */
app.all(
  '/projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId/*',
  (req: Request, res: Response) => {
    processParams(req, res);
    processMockRequest(req, res);
  }
);
/**
 * Health check for the service
 */
app.get('/healthz', (req: Request, res: Response) => {
  res.sendStatus(200);
  res.end();
});

app.listen(PORT, () =>
  console.log(`Server ready at: http://localhost:${PORT}`)
);
