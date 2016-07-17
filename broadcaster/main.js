/*jslint node: true, esversion: 6 */
'use strict';

var http = require('http');
var jsonBody = require('body/json');

var validator = require('./validator');
var broadcaster = require('./broadcaster');

// TODO:DLM: bring in from command line
let CREDENTIALS = '../credentials/msu-cave-f3ae939d1917.json';
let DATABASE_URL = 'https://msu-cave.firebaseio.com';
let PORT = 3000;

// TODO:DLM: the path is constructed from path concating the
// installation and the headset id
let path = 'test_data/yeah';
let firebaseRef = new broadcaster.FirebaseRef(CREDENTIALS, DATABASE_URL);
let firebaseBroadcaster = new broadcaster.Broadcaster(firebaseRef.database(), path);

function echo(output) {
  return function() { console.log(output); };
}

function finalizeResponse(request, response, statusCode, body) {
  response.statusCode = statusCode;
  response.write(JSON.stringify(body));
  console.log("Processed:" + request.url + " status:" + response.statusCode);
  response.end();
}

function broadcast(request, response, body) {
  firebaseBroadcaster.publish(body);
  finalizeResponse(request, response, 200, {});
}

function processValidBody(request, response, body) {
  let validData = validator.validate(body);
  if (! validData.success) {
    finalizeResponse(request, response, 400, validData.response);
  } else {
    broadcast(request, response, body);
  }
}

function handleRequest(request, response) {
  function processBody(err, body) {
    if (err) {
      let responseBody = {
        errorMessage: "Error processing body",
        developer: err
      };
      finalizeResponse(request, response, 500, responseBody);
    } else {
      processValidBody(request, response, body);
    }
  }
  jsonBody(request, response, processBody);
}

var server = http.createServer(handleRequest);
server.listen(PORT, echo("Server init"));
