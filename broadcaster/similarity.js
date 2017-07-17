/*jslint node: true, esversion: 6 */
'use strict';

const zip = (x, y) => x.map((xi, i) => [xi, y[i]]);

const correlationCoeff = (x, y) => {
  if (x.length != y.length) { throw "Lengths must be equal"; }

  const xBar = expectedValue(x);
  const yBar = expectedValue(y);

  const xDelta = x.map((xi) => xi - xBar);
  const yDelta = y.map((yi) => yi - yBar);

  const n = x.length;
  const num = zip(xDelta, yDelta)
    .reduce((sum, pi) => sum + (pi[0] * pi[1]), 0);

  const sigmaFromDelta = (delta) => {
    const exp_diff = delta
      .map((ai) => ai*ai)
      .reduce((sum, aiSq) => sum + aiSq, 0);
    return Math.sqrt(exp_diff);
  }

  const sigmaX = sigmaFromDelta(xDelta);
  const sigmaY = sigmaFromDelta(yDelta);

  return (sigmaX === 0 || sigmaY === 0) ? .5 : num / (sigmaX * sigmaY);
};

const expectedValue = (x) => {
  return x.reduce((sum, xi) => sum + xi) / x.length;
};

const align = (signal1, signal2) => {
  let alignment = [undefined, undefined];
  if (signal1.enoughSamples() && signal2.enoughSamples()) {
    const times = new Set([...signal1.times(), ...signal2.times()]);
    const orderedTimes = [...times].sort();
    const subSample = (signal) => orderedTimes.map((ti) => signal.eval(ti));
    alignment = [subSample(signal1), subSample(signal2)];
  }
  return alignment;
};

class Sample {
  constructor(time, value) {
    this.time = time;
    this.value = value;
  }

  isEqual(other) {
    return other && this.time === other.time && this.value === other.value;
  }
}

class Signal {
  constructor(windowSize) {
    this.samples = new Array();
    this.windowSize = windowSize;
  }

  times() {
    return this.samples.map((sample) => sample.time);
  }

  lastSample() {
    const tailIdx = this.samples.length - 1;
    return this.samples[tailIdx];
  }

  firstSample() {
    return this.samples[0];
  }

  enoughSamples() {
    return this.samples.length >= this.windowSize
  }

  eval(time) {
    let value;
    if (!this.enoughSamples()) {
      value = undefined;
    } else if (time <= this.firstSample().time) {
      value = this.firstSample().value;
    } else if (this.lastSample().time <= time) {
      value = this.lastSample().value;
    } else {
      const  index = this.samples.findIndex((sample) => time < sample.time);
      const [begin, end] = this.samples.slice(index-1, index+1);
      const t = (time-begin.time) / (end.time - begin.time);
      value = begin.value + t * (end.value - begin.value);
    }
    return value
  }

  addSample(time, value, headsetOn) {
    if (this.lastSample() && time < this.lastSample().time) {
      throw "Time must move forward";
    }

    if (headsetOn) {
      const sample = new Sample(time, value);
      this.samples.push(sample);
      if (this.samples.length > this.windowSize) {
        this.samples.shift();
      }
    } else {
      this.samples = [];
    }
  }
}

class SignalBank {
  constructor(localID, windowSize) {
    this.windowSize = windowSize;
    this.localID = localID;
    this.signals = new Map();
  }

  getSignal(key) {
    let signal = this.signals[key];
    if (!signal) {
      signal = new Signal(this.windowSize);
      this.signals[key] = signal;
    }
    return signal;
  }

  addSamples(samples) {
    for (let key in samples) {
      const rawData = samples[key].raw_data;
      const value = rawData.delta;
      const time = rawData.timestamp;
      const headsetOn = rawData.headsetOn;

      const target = this.getSignal(key);
      target.addSample(time, value, headsetOn);
    }
  }

  getActiveRemoteSignals() {
    const activeRemoteSignals = [];
    for (let key in this.signals) {
      const signal = this.getSignal(key)
      if (key !== this.localID && signal.enoughSamples()) {
        activeRemoteSignals.push(this.getSignal(key));
      }
    }
    return activeRemoteSignals;
  }

  similarity() {
    const localSignal = this.getSignal(this.localID);
    const remoteSignals = this.getActiveRemoteSignals();

    let sim = 0;
    if (localSignal.enoughSamples()) {
      const toAbsCorrelationCoeff = (signal) => {
        const [v1, v2] = align(localSignal, signal);
        return Math.abs(correlationCoeff(v1, v2));
      };

      sim = remoteSignals
        .map(toAbsCorrelationCoeff)
        .reduce((agg, corr) => agg + (corr / remoteSignals.length), 0);
    }
    return sim;
  }
}

module.exports.expectedValue = expectedValue;
module.exports.correlationCoeff = correlationCoeff;
module.exports.align = align;
module.exports.Signal = Signal;
module.exports.Sample = Sample;
module.exports.SignalBank = SignalBank;
