/*jslint node: true, esversion: 6 */
'use strict';

var yargs = require('yargs');
var osc = require('node-osc');
var fs = require('fs');

var broadcaster = require('./broadcaster');
var server = require('./server');

// setup the argument paring
var argv = yargs
  .usage('Usage: $0 [options]')
  // credentials
  .default('c', '../credentials/msu-cave-f3ae939d1917.json')
  .alias('c', 'credentials')
  .describe('c', 'Credential file for connecting to Firebase')
  // firebase server
  .default('f', 'https://msu-cave.firebaseio.com')
  .alias('f', 'firebase_url')
  .describe('f', 'Url of the firebase server')
  // port
  .default('p', 3000)
  .alias('p', 'port')
  .describe('p', 'Port on which to run the service')
  // installation id
  .demand('i')
  .alias('i', 'installation-id')
  .describe('i', 'The installation id')
  // headset id
  .demand('e')
  .alias('e', 'eeg-headset-id')
  .describe('e', 'The eeg headset id')
  //osc clients
  .default('o', [])
  .alias('o', 'osc-servers')
  .describe('o', 'the osc servers to send osc data')
  .array('o')
  // help
  .help('h')
  .alias('h', 'help')
  .argv;

// create and initialize the broadcaster object.  The broadcaster handles
// much of the heavy lifting for communicating updates on the node to
// the firebase serer
let db = new broadcaster.firebaseDB(argv.credentials, argv.firebase_url);
let firebaseBroadcaster = new broadcaster.Broadcaster(
  db, argv.installationId, argv.eegHeadsetId
);

// initialize so that every time remote data is updated the onRemoteData
// method is called.  This should hook into the covariance calculator either
// by sending a message to the covariance sevice or calling directly.
function onRemoteData(snapshot) {
  // NOTE:DLM: whoever is doing the covariance, this is where you can hook
  // in your code or call your service.  If you are going to call your
  // service, I recommend using the requests module
  // https://github.com/request/request
  console.log(snapshot.val());
}
firebaseBroadcaster.subscribe(onRemoteData);

// setup the osc clients
var clients = argv.oscServers.map(function(rawClientAddress) {
  let clientAddress = rawClientAddress.split(':');
  return new osc.Client(clientAddress[0], clientAddress[1]);
});

// setup the server so that everything it receives some new data it is
// published to the remote data server.
function onLocalData(body) {
  firebaseBroadcaster.publish(body);

  // Send an OSC message
  var data = [
    body.timestamp,
    body.delta,
    body.hiAlpha,
    body.hiBeta,
    body.loAlpha,
    body.loBeta,
    body.loGamma,
    body.midGamma,
    body.theta
  ];
  console.log(data);


  clients.forEach(function(client) { client.send("/eeg", data); });
}
var server = new server.Server(argv.port, onLocalData);
server.start();
