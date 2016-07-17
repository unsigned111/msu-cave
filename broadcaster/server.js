/*jslint node: true, esversion: 6 */
'use strict';

var http = require('http');
var jsonBody = require('body/json');

var validator = require('./validator');

class Server {

  _finalizeResponse(request, response, statusCode, body) {
    response.statusCode = statusCode;
    response.write(JSON.stringify(body));
    console.log("Processed:" + request.url + " status:" + response.statusCode);
    response.end();
  }

  _handleServerError(request, response, message, error) {
    let responseBody = {
      errorMessage: message,
      developer: error
    };
    this._finalizeResponse(request, response, 500, responseBody);
  }

  _processValidBody(request, response, body) {
    let validData = validator.validate(body);
    if (! validData.success) {
      this._finalizeResponse(request, response, 400, validData.response);
    } else {
      this.onLocalData(body);
      this._finalizeResponse(request, response, 200, {});
    }
  }

  _makeHandleRequest() {
    var that = this;
    function handleRequest(request, response) {
      function processBody(err, body) {
        if (err) {
          that._handleServerError(request, response, "Error processing body", err);
        } else {
          that._processValidBody(request, response, body);
        }
      }
      jsonBody(request, response, processBody);
    }
    return handleRequest;
  }

  /**
   * create an instance of a broadcast server.
   * @param port the port on which to run the server
   * @param onLocalData (body) -> void that will run whenever valid data
   *   is received.  Often a good choice here is calling a broadcaster.
   */
  constructor(port, onLocalData) {
    this.port = port;
    let handleRequest = this._makeHandleRequest();
    this.server = http.createServer(handleRequest);
    this.onLocalData = onLocalData;
  }

  start() {
    this.server.listen(this.port, () => console.log("Server init"));
  }
}

module.exports.Server = Server;
