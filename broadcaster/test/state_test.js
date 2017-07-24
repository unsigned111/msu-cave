/*jslint node: true, esversion: 6 */

const assert = require('chai').assert;

const state = require('../state');

suite('state', function() {
  suite('OnOffModel', function() {
    suite('constructor', function() {
      test('it initializes', function() {
        const m = new state.OnOffModel(.5, 3);
        assert.equal(.5, m.threashold);
        assert.deepEqual([255, 255, 255], m.samples);
      });
    });

    suite('addSample', function() {
      test('buffer is filled from back', function() {
        const m = new state.OnOffModel(.5, 3);
        m.addSample(8);
        assert.deepEqual([255, 255, 8], m.samples);
      });
    });


    suite('isOn', function() {
      test('it is on for low values', function() {
        const m = new state.OnOffModel(.5, 2);
        m.samples = [4, 10];
        assert.equal(true, m.isOn());
      });

      test('it is off for high values', function() {
        const m = new state.OnOffModel(.5, 2);
        m.samples = [225, 250];
        assert.equal(false, m.isOn());
      });
    });
  });

  suite('State', function() {
    const data = {
      attention: 1,
      delta: 2,
      hiAlpha: 3,
      hiBeta: 4,
      loAlpha: 5,
      loBeta: 6,
      loGamma: 7,
      meditation: 8,
      midGamma: 9,
      signal: 10,
      theta: 11,
      timestamp: 12
    };

    const makeState = () => {
      const m = new state.OnOffModel(.5, 1);
      return new state.State(m);
    };

    const makeAndFillState = () => {
      const s = makeState();
      s.addData(data);
      return s
    };

    suite('#addData', function() {
      test('it updates the state', function() {
        const s = makeState();
        s.addData(data);

        assert.equal(1, s.attention);
        assert.equal(2, s.delta);
        assert.equal(3, s.hiAlpha);
        assert.equal(4, s.hiBeta);
        assert.equal(5, s.loAlpha);
        assert.equal(6, s.loBeta);
        assert.equal(7, s.loGamma);
        assert.equal(8, s.meditation);
        assert.equal(9, s.midGamma);
        assert.equal(10, s.signal);
        assert.equal(11, s.theta);
        assert.equal(12, s.timestamp);
        assert.deepEqual([2], s.onOffModel.samples);
      });
    });

    suite('#toOscEeg', function() {
      const s = makeAndFillState();
      assert.deepEqual([12, 2, 3, 4, 5, 6, 7, 9, 11], s.toOscEeg());
    });

    suite('#toOscOnOff', function() {
      const s = makeAndFillState();
      assert.deepEqual([1], s.toOscOnOff());
    });
  });
});
