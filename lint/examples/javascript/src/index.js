// This loads definitions of the structs in our protocol buffer models.
const Style = require('../pbjs-genfiles/proto').google.cloud.apigeeregistry.v1.style;

async function main() {

  // The plugin expects a request message on stdin.
  process.stdin.on("data", data => {
  
    // The request should be an encoded LinterRequest.
    // Decode it.
    const request = Style.LinterRequest.decode(
      data
    );
  
    // Log it.
    process.stderr.write(JSON.stringify(request) + "\n");

    // We're going to make up some fake results.
    // Here we one sample violation of each rule that was listed in the request.
    problems = request.ruleIds.map(x => Style.LintProblem.fromObject({
      "ruleId": x, 
      "message": "it is violated",
      "suggestion": "keep API-ing!",
      "location": {
        "startPosition": {
          "lineNumber": 2, 
          "columnNumber": 3
        },
        "endPosition": {
          "lineNumber": 4,
          "columnNumber": 5,
        }
      }
    }))

    // The request included a directory that contains the spec to lint,
    // but we would have to look there to get its file name.
    // Assume for now that we are linting a file named "swagger.yaml".
    file = Style.LintFile.fromObject({
      "filePath": request.specDirectory + "/swagger.yaml",
      "problems": problems
    });

    // Put this all together in a response message.
    response = Style.LinterResponse.fromObject({
      "lint": {
        "name": "registry-lint.js",
        "files": [file]
      }
    });

    // Log the response message.
    process.stderr.write(JSON.stringify(response) + "\n");

    // Encode the response and write it to stdout.
    const responseBuffer = Style.LinterResponse.encode(
      response
    ).finish();
    process.stdout.write(responseBuffer)

    // We're done.
    process.exit();
  })

  // We just need this to wait until the handler above finishes.
  await sleep(1000);
}

main().catch(err => {
  console.error(err);
});

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}
