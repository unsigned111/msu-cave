/*jslint node: true, esversion: 6 */

const assert = require('chai').assert;

const similarity = require('../similarity');

suite('simiarity', function () {
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
    test('it returns undefined when not enough data')
    test('it returns aligned samples when enough data')
  });

  suite('Signal', function() {
    const makeSample = (time, value) => { return { time, value } };

    suite('#lastSample', function() {
      test('returns undefined when no samples', function() {
        const signal = new similarity.Signal(2);
        assert.equal(undefined, signal.lastSample());
      });

      test('it returns value when present', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(1, 2);
        assert.deepEqual(makeSample(1,2), signal.lastSample());
      });
    });

    suite('#times', function() {
      test('it returns times', function() {
        const signal = new similarity.Signal(2);
        signal.addSample(1, 2);
        signal.addSample(3, 5);
        assert.deepEqual([1,3], signal.times());
      });

    });

    suite('#addSample', function() {
      test('maintins circular list', function() {
        const signal = new similarity.Signal(3);

        signal.addSample(1, 2);
        signal.addSample(2, 4);
        signal.addSample(3, 8);
        assert.deepEqual([
          makeSample(1, 2),
          makeSample(2, 4),
          makeSample(3, 8),
        ], signal.samples)

        signal.addSample(4, 16);
        assert.deepEqual([
          makeSample(2, 4),
          makeSample(3, 8),
          makeSample(4, 16),
        ], signal.samples)
      });
    });

    suite('#eval', function() {
      describe('not enough samples', function() {
        test('it is undefined with too few samples', function() {
          const signal = new similarity.Signal(2);
          signal.addSample(1, 2);
          assert.equal(undefined, signal.eval(6));
        });
      });

      describe('enough samples', function() {
        const makeSignal = () => {
          const signal = new similarity.Signal(3);
          signal.addSample(2, 4);
          signal.addSample(3, 8);
          signal.addSample(4, 16);
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
