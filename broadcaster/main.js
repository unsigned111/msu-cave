/*jslint node: true, esversion: 6 */
'use strict';

var yargs = require('yargs');
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

// setup the osc clients. The osc broadcaser hands the heavy lifting
// of sending the data over to the OSC clients
const clients = argv.oscServers.map(broadcaster.oscClient);
const oscBroadcaster = new broadcaster.OSCBroadcaster(clients);

// create and initialize the broadcaster object.  The broadcaster handles
// much of the heavy lifting for communicating updates on the node to
// the firebase serer
const db = new broadcaster.firebaseDB(argv.credentials, argv.firebase_url);
const firebaseBroadcaster = new broadcaster.FirebaseBroadcaster(
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
  // console.log(snapshot.val());
}
firebaseBroadcaster.subscribe(onRemoteData);

// setup the server so that everything it receives some new data it is
// published to the remote data server.
function onLocalData(body) {
  console.log(body)
  firebaseBroadcaster.publish(body);
  oscBroadcaster.publish(body);
}
const webServer = new server.Server(argv.port, onLocalData);
webServer.start();
