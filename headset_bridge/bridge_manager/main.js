#!/usr/bin/env node

/*jslint node: true, esversion: 6 */
'use strict';


const yargs = require('yargs');
const osc = require('node-osc');

const monitor = require('./monitor');
const service = require('./service');

// setup the argument paring
const argv = yargs
  .usage('Usage: $0 [options]')
  // port
  .default('p', 7772)
  .alias('p', 'headset-bridge-port')
  .describe('p', 'Port for listenting for headset bridge')
  // headset hardware address
  .demand('a')
  .alias('a', 'target-address')
  .describe('a', 'Target headset hardware address.')
  // timeout
  .default('t', 20)
  .alias('t', 'timeout')
  .describe('t', 'timeout length for how long to wait for a restart.')
  // help
  .help('h')
  .alias('h', 'help')
  .argv;


const headset = {
  name: "headset",
  start: `cd ../..; pwd; ./headset_bridge/attach ${argv.targetAddress}`,
  stop: "cd ../..; ./headset_bridge/unattach",
};
const headsetService = new service.Service(headset);
const headsetMonitor = new monitor.Monitor(headsetService, argv.timeout*1000);

const oscServer = new osc.Server(argv.headsetBridgePort);
oscServer.on("message", function (msg, rinfo) {
  const oscAddress = msg[0];
  if (oscAddress === '/eeg') {
    headsetMonitor.update();
  }
});

setInterval(() => { headsetMonitor.restartIfTimedOut() }, 1000);
