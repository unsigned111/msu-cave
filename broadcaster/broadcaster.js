/*jslint node: true, esversion: 6 */
'use strict';

var firebase = require('firebase');

class FirebaseRef {
  constructor(credentials, url) {
    this.app = firebase.initializeApp({
      serviceAccount: credentials,
      databaseURL: url,
    });
  }

  database() {
    return firebase.database();
  }
}

class Broadcaster {
  constructor(database, path) {
    this.db = database;
    this.path = path;
    this.ref = this.db.ref(path);
  }

  publish(data) {
    this.ref.set({ cur_data: data });
  }

  subscribe(callback) {
    this.ref.on('value', callback);
  }
}

module.exports.FirebaseRef = FirebaseRef;
module.exports.Broadcaster = Broadcaster;

