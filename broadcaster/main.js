/*jslint node: true, esversion: 6 */
'use strict';

var broadcaster = require('./broadcaster');
var server = require('./server');

// TODO:DLM: bring in from command line
let CREDENTIALS = '../credentials/msu-cave-f3ae939d1917.json';
let DATABASE_URL = 'https://msu-cave.firebaseio.com';
let INSTALLATION_ID = 'installation-id';
let HEADSET_ID = 'headset-id';
let PORT = 3000;

let db = new broadcaster.firebaseDB(CREDENTIALS, DATABASE_URL);
let firebaseBroadcaster = new broadcaster.Broadcaster(db, INSTALLATION_ID, HEADSET_ID);

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

// setup the server so that everyting it receives some new data it is
// published to the remote data.
function onLocalData(body) {
  firebaseBroadcaster.publish(body);
}
var server = new server.Server(PORT, onLocalData);
server.start();
