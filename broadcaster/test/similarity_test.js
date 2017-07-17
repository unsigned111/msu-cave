/*jslint node: true, esversion: 6 */

const assert = require('chai').assert;

const similarity = require('../similarity');

suite('similarity', function () {
  suite('#expectedValue', function() {
    test('it returns the expected value', function() {
      const result = similarity.expectedValue([1,2,3,4,5]);
      assert.equal(((5*6)/2) / 5, result)
    });
  })

  describe('#correlationCoeff', function() {
    test('it returns the correlation coeff', function() {
      const x = [1.0, 5.5, 7.8, 4.2, -2.7, -5.4, 8.9];
      const y = [0.1, 1.5, 0.8, -4.2, 2.7, -9.4, -1.9];
      const result = similarity.correlationCoeff(x, y);
      assert.approximately(0.3893218, result, 0.0001);
    });

    test('it does not freak out for zero sigma', function() {
      const x = [1.0, 1.0];
      const y = [0.1, 0.78];
      const result = similarity.correlationCoeff(x, y);
      assert.isNotNaN(result);
    });
  });

  describe('#align', function() {
    test('it returns aligned samples when enough data', function() {
      const signals = [new similarity.Signal(2), new similarity.Signal(3)];
      signals[0].addSample(1, 4, true);
      signals[0].addSample(4, 1, true);
      signals[1].addSample(2, 3, true);
      signals[1].addSample(4, 6, true);
      signals[1].addSample(8, 12, true);
      const [v1, v2] = similarity.align(signals[0], signals[1])

      assert.deepEqual([4, 3, 1, 1], v1)
      assert.deepEqual([3, 3, 6, 12], v2)
    });

    test('it returns undefined when not enough data', function() {
      const signals = [new similarity.Signal(2), new similarity.Signal(3)];
      const [v1, v2] = similarity.align(signals[0], signals[1])
      assert.equal(undefined, v1)
      assert.equal(undefined, v2)
    });
  });

  suite('Signal', function() {
    const makeSample = (time, value) => new similarity.Sample(time, value);

    suite('#enounghSamples', function() {
      test('it returns true with enouth samples', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(1, 2, true);
        signal.addSample(2, 4, true);
        assert.equal(true, signal.enoughSamples());
      });

      test('it returns false with enouth samples', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(2, 4, true);
        assert.equal(false, signal.enoughSamples());
      });
    });

    suite('#lastSample', function() {
      test('it returns undefined when no samples', function() {
        const signal = new similarity.Signal(2);
        assert.equal(undefined, signal.lastSample());
      });

      test('it returns value when present', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(1, 2, true);
        assert.deepEqual(makeSample(1,2), signal.lastSample());
      });
    });

    suite('#times', function() {
      test('it returns times', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(1, 2, true);
        signal.addSample(3, 5, true);
        assert.deepEqual([1,3], signal.times());
      });

    });

    suite('#addSample', function() {
      test('maintins circular list', function() {
        const signal = new similarity.Signal(3);

        signal.addSample(1, 2, true);
        signal.addSample(2, 4, true);
        signal.addSample(3, 8, true);
        assert.deepEqual([
          makeSample(1, 2),
          makeSample(2, 4),
          makeSample(3, 8),
        ], signal.samples)

        signal.addSample(4, 16, true);
        assert.deepEqual([
          makeSample(2, 4),
          makeSample(3, 8),
          makeSample(4, 16),
        ], signal.samples)
      });

      test('resets on a headset off', function() {
        const signal = new similarity.Signal(3);

        signal.addSample(1, 2, true);
        assert.deepEqual([makeSample(1, 2)], signal.samples)

        signal.addSample(4, 0, false);
        assert.deepEqual([], signal.samples)
      });
    });

    suite('#eval', function() {
      describe('not enough samples', function() {
        test('it is undefined with too few samples', function() {
          const signal = new similarity.Signal(2);
          signal.addSample(1, 2, true);
          assert.equal(undefined, signal.eval(6));
        });
      });

      describe('enough samples', function() {
        const makeSignal = () => {
          const signal = new similarity.Signal(3);
          signal.addSample(2, 4, true);
          signal.addSample(3, 8, true);
          signal.addSample(4, 16, true);
          return signal;
        }

        test('it is a value within an interval', function() {
          const signal = makeSignal();
          assert.equal(10, signal.eval(3.25));
          assert.equal(12, signal.eval(3.5));
          assert.equal(14, signal.eval(3.75));
        });

        test('it is a value at a larger time', function() {
          const signal = makeSignal();
          assert.equal(16, signal.eval(5));
        });

        test('it is a value at a smaller time time', function() {
          const signal = makeSignal();
          assert.equal(4, signal.eval(0));
        });
      });
    });
  });

  suite('SignalBank', function() {
    const makeRawData = (val, headsetOn, timeDelta=0) => {
      return {
        raw_data: {
          delta: val, headsetOn, timestamp: Date.now()+timeDelta,
        },
      };
    };

    const makeNewSample = () => {
      return {
        e: makeRawData(1578003, true),
        f: makeRawData(1578003, true),
      };
    };

    suite('#constructor', function() {
      it('sets the local id', function() {
        const bank = new similarity.SignalBank('myLocalId', 2);
        assert.equal('myLocalId', bank.localID);
      });

      it('sets the window size', function() {
        const bank = new similarity.SignalBank('myLocalId', 2);
        assert.equal(2, bank.windowSize);
      });
    });

    suite('#getSignal', function() {
      const localID = 'localID';

      test('it creats a signal if there is not one', function() {
        const bank = new similarity.SignalBank(localID, 2);
        bank.getSignal(localID);
        assert(bank.signals[localID])
      });

      test('it returns signal if there is one', function() {
        const bank = new similarity.SignalBank(localID, 2);
        const signal = bank.getSignal(localID);
        assert.equal(signal, bank.getSignal(localID));
      });
    });

    suite('#addSamples', function() {
      const localID = 'localID';

      test('it adds a sample', function() {
        const bank = new similarity.SignalBank(localID, 2);
        const sample = { localID: makeRawData(10, true) };

        bank.addSamples(sample);
        assert.deepEqual(
          [{ time: sample[localID].raw_data.timestamp, value: 10 }],
          bank.signals[localID].samples
        )
      });
    });

    suite('#getActiveRemoteSignals', function() {
      const bank = new similarity.SignalBank('local', 1);

      bank.addSamples({
        local: makeRawData(10, true),
        remote: makeRawData(11, true),
      });
      const remoteSignal = bank.getSignal('remote');
      bank.getSignal('remoteInactive');

      assert.deepEqual([remoteSignal], bank.getActiveRemoteSignals());
    })

    suite('#similarity', function() {
      test('when local is inactive it returns undefined');
      test('when the headset is active and no others around it returns undefined')

      test('it returns value when active and others around ', function() {
        const bank = new similarity.SignalBank('local', 4);

        bank.addSamples({
          local: makeRawData(15151, true),
          remote1: makeRawData(30567, true),
          remote2: makeRawData(17206, true),
          remote3: makeRawData(1262343, false),
        });
        bank.addSamples({
          local: makeRawData(152782, true, 5),
          remote1: makeRawData(814364, true, 5),
          remote2: makeRawData(1252675, true, 5),
          remote3: makeRawData(632536, false, 5),
        });
        bank.addSamples({
          local: makeRawData(67037, true, 10),
          remote1: makeRawData(1170322, true, 10),
          remote2: makeRawData(1243019, true, 10),
        });
        bank.addSamples({
          local: makeRawData(630661, true, 15),
          remote1: makeRawData(1187917, true, 15),
          remote2: makeRawData(287963, true, 15),
        });

        const sim = bank.similarity();
        assert.approximately(0.420775, sim, 0.0001);
      });
    });
  });
});
