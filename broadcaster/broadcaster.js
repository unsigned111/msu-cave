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

  publishHeadset(data) {
    // Send an eeg OSC message
    const eegData = [
      data.timestamp,
      data.delta,
      data.hiAlpha,
      data.hiBeta,
      data.loAlpha,
      data.loBeta,
      data.loGamma,
      data.midGamma,
      data.theta,
    ];
    this.publishToAll("/eeg", eegData);

    // Send an on/off OSC message
    const onOffData = [
      data.headsetOn ? 1 : 0,
    ]
    this.publishToAll("/onoff", onOffData);
  }
}

module.exports.firebaseDB = firebaseDB;
module.exports.FirebaseBroadcaster = FirebaseBroadcaster;
module.exports.oscClient = oscClient;
module.exports.OSCBroadcaster = OSCBroadcaster;
