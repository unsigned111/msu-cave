/*jslint node: true, esversion: 6 */


const log = (msg) => console.log(msg);

class Monitor {
  constructor(service, timeout) {
    this.service = service;
    this.timeout = timeout;
    this.lastUpdate = 0;
    this.update();
  }

  now() {
    return (new Date()).getTime();
  }

  update() {
    this.lastUpdate = this.now();
  }

  restartService() {
    this.service.stop();
    this.service.start();
    this.update();
  }

  deltaT() {
    return this.now() - this.lastUpdate;
  }

  hasTimedOut() {
    return this.deltaT() > this.timeout;
  }

  restartIfTimedOut() {
    const restart = this.hasTimedOut();
    if (restart) {
      log(`restarting "${this.service.name}" no response for ${this.deltaT() / 1000}`);
      this.restartService();
    }
    return restart;
  }
}

module.exports.Monitor = Monitor;
