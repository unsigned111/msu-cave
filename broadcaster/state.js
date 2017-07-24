/*jslint node: true, esversion: 6 */
'use strict';

const MAX_SAMPLE = 255;

class OnOffModel {
  constructor(threashold, windowSize) {
    this.threashold = threashold;
    this.samples = new Array(windowSize).fill(MAX_SAMPLE);
  }

  addSample(sample) {
    this.samples.push(sample);
    this.samples.shift();
  }

  isOn() {
    const agg = this.samples.reduce((agg, si) => agg + si / MAX_SAMPLE, 0)
    const average = agg / this.samples.length;
    return  average < this.threashold;
  }
}


class State {
  constructor(onOffModel) {
    this.attention = 0;
    this.delta = 0;
    this.hiAlpha = 0;
    this.hiBeta = 0;
    this.loAlpha = 0;
    this.loBeta = 0;
    this.loGamma = 0;
    this.meditation = 0;
    this.midGamma = 0;
    this.signal = 0;
    this.theta = 0;
    this.timestamp = 0;
    this.onOffModel = onOffModel;
  }

  addData(data) {
    this.attention = data.attention;
    this.delta = data.delta;
    this.hiAlpha = data.hiAlpha;
    this.hiBeta = data.hiBeta;
    this.loAlpha = data.loAlpha;
    this.loBeta = data.loBeta;
    this.loGamma = data.loGamma;
    this.meditation = data.meditation;
    this.midGamma = data.midGamma;
    this.signal = data.signal;
    this.theta = data.theta;
    this.timestamp = data.timestamp;
    this.onOffModel.addSample(data.signal)
  }

  toOscEeg() {
    return [
      this.timestamp,
      this.delta,
      this.hiAlpha,
      this.hiBeta,
      this.loAlpha,
      this.loBeta,
      this.loGamma,
      this.midGamma,
      this.theta,
    ];
  }

  toOscOnOff() {
    return [
      this.onOffModel.isOn() ? 1 : 0,
    ];
  }
}

module.exports.OnOffModel = OnOffModel;
module.exports.State = State;
