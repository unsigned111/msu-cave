/*jslint node: true, esversion: 6 */
'use strict';

var http = require('http');
var jsonBody = require('body/json');

var validator = require('./validator');
var broadcaster = require('./broadcaster');

// TODO:DLM: bring in from command line
let CREDENTIALS = '../credentials/msu-cave-f3ae939d1917.json';
let DATABASE_URL = 'https://msu-cave.firebaseio.com';
let INSTALLATION_ID = 'installation-id';
let HEADSET_ID = 'headset-id';
let PORT = 3000;

let db = new broadcaster.firebaseDB(CREDENTIALS, DATABASE_URL);
let firebaseBroadcaster = new broadcaster.Broadcaster(db, INSTALLATION_ID, HEADSET_ID);


function onNewData(snapshot) {
  console.log(snapshot.val());
}
firebaseBroadcaster.subscribe(onNewData);

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
server.listen(PORT, () => console.log("Server init"));
