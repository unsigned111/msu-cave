/*jslint node: true, esversion: 6 */
'use strict';

var firebase = require('firebase');

function firebaseDB(credentials, url) {
  firebase.initializeApp({
    serviceAccount: credentials,
    databaseURL: url,
  });
  return firebase.database();
}

class Broadcaster {
  constructor(database, installationID, headsetID) {
    this.db = database;
    this.installationID = installationID;
    this.headsetID = headsetID;
    this._ref = this.db.ref(this.latestPath);
  }

  publish(data) {
    let payload = {};
    payload[this.headsetID] = data;
    this._ref.set(payload);
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

module.exports.firebaseDB = firebaseDB;
module.exports.Broadcaster = Broadcaster;

