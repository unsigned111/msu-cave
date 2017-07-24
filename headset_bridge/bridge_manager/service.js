/*jslint node: true, esversion: 6 */

const exec = require('child_process').execSync;

class Service {
  constructor(config) {
    this.name = config.name;
    this.startCMD = config.start;
    this.stopCMD = config.stop;
  }

  runCmd(cmd) {
    exec(cmd, {stdio: [0,1,2]});
  }

  stop() {
    this.runCmd(this.stopCMD);
  }

  start() {
    this.runCmd(this.startCMD);
  }
}

module.exports.Service = Service;
