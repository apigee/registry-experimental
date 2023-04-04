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
import {Request, Response} from 'express';
import * as cors from 'cors';
import * as fs from 'fs';
import * as path from 'path';

import {RegistryClient} from '@google-cloud/apigee-registry';
import {ClientOptions} from 'google-gax';
import {credentials} from '@grpc/grpc-js';
import {makeExecutableSchema} from '@graphql-tools/schema';
import {addMocksToSchema} from '@graphql-tools/mock';
import {graphql} from 'graphql';

const app = express();

app.use(cors());
app.use(express.json());

const PORT = process.env.PORT || 3000;

const client_options = <ClientOptions>{};

if (process.env.REGISTRY_INSECURE) {
  if (process.env.REGISTRY_INSECURE === '1') {
    client_options.sslCreds = credentials.createInsecure();
  }
}

if (process.env.REGISTRY_ADDRESS) {
  const items = process.env.REGISTRY_ADDRESS.split(':');
  client_options.apiEndpoint = items[0];
  client_options.port = items.length >= 1 ? parseInt(items[1]) : 443;
}

const client = new RegistryClient(client_options);

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
 * and pass the spec to the GraphQL Mocking Library.
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
  //Used cached version of the file if it exists.
  const spec_local_path = '/tmp/' + spec_url;
  fs.readFile(spec_local_path, (err, contents) => {
    if (!err && contents) {
      executeMockRequest(contents.toString(), req, res);
    } else {
      client.getApiSpecContents(
        {
          name: spec_url,
        },
        (err, response) => {
          if (err) {
            return _sendError(err, res);
          } else {
            if (response && response.data) {
              const specString = response.data.toString();
              fs.mkdirSync(path.dirname(spec_local_path), {recursive: true});
              fs.writeFileSync(spec_local_path, specString);
              executeMockRequest(specString, req, res);
            } else {
              return _sendError(new Error('Spec not found'), res);
            }
          }
        }
      );
    }
  });
}

/**
 * Process request using graphql-tools mock library
 *
 * @param specString
 * @param req
 * @param res
 */
function executeMockRequest(
  specString: string,
  req: express.Request,
  res: express.Response
) {
  const schema = makeExecutableSchema({typeDefs: specString});

  const mocks = {
    Date: () => new Date(),
  };
  const schemaWithMocks = addMocksToSchema({schema, mocks});
  graphql({
    schema: schemaWithMocks,
    source: req.body.query,
  })
    .then(result => {
      res.send(result).end();
    })
    .catch(err => {
      console.error(err);
      res.sendStatus(500).end();
    });
}

function processParams(req: Request, res: Response) {
  for (const [key, value] of Object.entries(req.params)) {
    res.locals[key] = value;
  }
}

/**
 * Mock request endpoint for API Deployments
 *
 * This handler will lookup the latest spec revision from the deployment and
 * forward the request to /projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId/*
 * handler
 */
app.all(
  '/projects/:projectId/locations/:locationId/apis/:apiId/deployments/:deploymentId',
  (req, res, next) => {
    processParams(req, res);
    const url =
      'projects/' +
      res.locals.projectId +
      '/locations/' +
      res.locals.locationId +
      '/apis/' +
      res.locals.apiId +
      '/deployments/' +
      res.locals.deploymentId;
    client.getApiDeployment(
      {
        name: url,
      },
      (err, response) => {
        if (response && response.apiSpecRevision) {
          /**
           * Pass on the handling of this request to the next route
           * /projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId/*
           */
          req.url = '/' + response.apiSpecRevision;
          next();
        } else {
          if (!err) {
            err = new Error('Error processing deployment[' + url + ']');
          }
          _sendError(err, res);
        }
      }
    );
  }
);

/**
 * Mock request endpoint with registry spec in request Path
 */
app.all(
  '/projects/:projectId/locations/:locationId/apis/:apiId/versions/:versionId/specs/:specId',
  (req: Request, res: Response) => {
    processParams(req, res);
    processMockRequest(req, res);
  }
);
/**
 * Health check for the service
 */
app.get('/', (req: Request, res: Response) => {
  res.sendStatus(200);
  res.end();
});
app.listen(PORT, () =>
  console.log(`Server ready at: http://localhost:${PORT}`)
);
