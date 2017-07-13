/*jslint node: true, esversion: 6 */
'use strict';

const zip = (x, y) => x.map((xi, i) => [xi, y[i]]);

const covariance = (x, y) => {
  if (x.length != y.length) { throw "Lengths must be equal"; }

  const xBar = expectedValue(x);
  const yBar = expectedValue(y);

  const xDelta = x.map((xi) => xi - xBar);
  const yDelta = y.map((yi) => yi - yBar);

  const n = x.length;
  const num = zip(xDelta, yDelta)
    .reduce((sum, pi) => sum + (pi[0] * pi[1]), 0);

  return num / (n-1);
};

const expectedValue = (x) => {
  return x.reduce((sum, xi) => sum + xi) / x.length;
};

class Signal {
  constructor(windowSize) {
    this.samples = new Array();
    this.windowSize = windowSize;
  }

  lastSample() {
    const tailIdx = this.samples.length - 1;
    return this.samples[tailIdx];
  }

  firstSample() {
    return this.samples[0];
  }

  eval(time) {
    let value;
    if (this.samples.length < this.windowSize) {
      value = undefined;
    } else if (time <= this.firstSample().time) {
      value = this.firstSample().value;
    } else if (this.lastSample().time < time) {
      value = this.lastSample().value;
    } else {
      const  index = this.samples.findIndex((sample) => time < sample.time);
      const [begin, end] = this.samples.slice(index-1, index+1);
      const t = (time-begin.time) / (end.time - begin.time);
      value = begin.value + t * (end.value - begin.value);
    }
    return value
  }

  addSample(time, value) {
    if (this.lastSample() && time < this.lastSample().time) {
      throw "Time must move forward";
    }

    const sample = this.makeSample(time, value);
    this.samples.push(sample);
    if (this.samples.length > this.windowSize) {
      this.samples.shift();
    }
  }

  makeSample(time, value) {
    return { time, value };
  }
}

module.exports.expectedValue = expectedValue;
module.exports.covariance = covariance;
module.exports.Signal = Signal;
