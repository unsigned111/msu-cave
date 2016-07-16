/*jslint node: true, esversion: 6 */
'use strict';

var firebase = require('firebase');

// TODO:DLM: bring in from command line
let CREDENTIALS = '../credentials/msu-cave-f3ae939d1917.json';
let DATABASE_URL = 'https://msu-cave.firebaseio.com';

var app = firebase.initializeApp({
    serviceAccount: CREDENTIALS,
    databaseURL: DATABASE_URL,
});

let path = 'test_data/yeah';
var db = firebase.database();
var ref = db.ref(path);
ref.set({ test_data: 'some data' });

ref.on('value', function(snapshot) {
  console.log(snapshot.val());
});

setTimeout(function() {
    ref.set({ test_data: 'some new data' });
}, 1000);
