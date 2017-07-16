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

  describe('#covariaiance', function() {
    test('it returns the covariance', function() {
      const x = [1.0, 5.5, 7.8, 4.2, -2.7, -5.4, 8.9];
      const y = [0.1, 1.5, 0.8, -4.2, 2.7, -9.4, -1.9];
      const result = similarity.covariance(x, y);
      assert(8.697381-result < 0.0001);
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
});
