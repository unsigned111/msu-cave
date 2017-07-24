/*jslint node: true, esversion: 6 */
'use strict';

var firebase = require('firebase');
var osc = require('node-osc');

function firebaseDB(credentials, url) {
  firebase.initializeApp({
    serviceAccount: credentials,
    databaseURL: url,
  });
  return firebase.database();
}

class FirebaseBroadcaster {
  constructor(database, installationID, headsetID) {
    this.db = database;
    this.installationID = installationID;
    this.headsetID = headsetID;
    this._ref = this.db.ref(this.latestPath);
  }

  publish(data) {
    let payload = {};
    payload = {
      raw_data: data,
      headsetOn: data.onOffModel.isOn(),
      timestamp: {
        server: firebase.database.ServerValue.TIMESTAMP,
        node: (new Date()).getTime()
      }
    };
    this._ref.child(this.headsetID).set(payload);
  }

  subscribe(callback) {
    this._ref.on('value', callback);
  }

  get installationPath() {
    return 'installations/' + this.installationID;
  }

  get latestPath() {
    return this.installationPath + '/latest';
  }
}


function oscClient(rawClientAddress) {
  let clientAddress = rawClientAddress.split(':');
  return new osc.Client(clientAddress[0], clientAddress[1]);
}

class OSCBroadcaster {
  constructor(clients) {
    this.clients = clients;
  }

  publishToAll(channel, data) {
    this.clients.forEach((client) => client.send(channel, data));
  }

  publishHeadset(state) {
    // Send an eeg OSC message
    const eegData = state.toOscEeg()
    this.publishToAll("/eeg", eegData);

    // Send an on/off OSC message
    const onOffData = state.toOscOnOff()
    this.publishToAll("/onoff", onOffData);
  }

  publishSimilarity(similarity) {
    this.publishToAll("/similarity", [similarity]);
  }
}

module.exports.firebaseDB = firebaseDB;
module.exports.FirebaseBroadcaster = FirebaseBroadcaster;
module.exports.oscClient = oscClient;
module.exports.OSCBroadcaster = OSCBroadcaster;
