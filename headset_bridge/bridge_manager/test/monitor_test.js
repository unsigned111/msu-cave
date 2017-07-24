/*jslint node: true, esversion: 6 */

const assert = require('chai').assert;
const sinon = require('sinon');

const monitor = require('../monitor');

suite('monitor', function() {
  suite('Monitor', function() {
    const makeService = () => {
      return { start: () => undefined, stop: () => undefined, name: "" };
    };

    suite('constructor', function() {
      test('initalized', function() {
        const service = makeService();
        const m = new monitor.Monitor(service, 5000);

        assert.equal(5000, m.timeout);
        assert.equal(service, m.service);
        assert.isBelow(0, m.lastUpdate);
      });
    });

    suite('#update', function() {
      test('updates to the current time', function() {
        const service = makeService();
        const m = new monitor.Monitor(service, 5000);
        const expectedTime = 123;
        sinon.stub(m, "now").callsFake(() => expectedTime);
        m.update();
        assert.equal(expectedTime, m.lastUpdate);
      });
    });

    suite('#restartService', function() {
      test('stops and starts the service', function() {
        const service = makeService();
        const mock = sinon.mock(service);

        mock.expects('stop');
        mock.expects('start');

        const m = new monitor.Monitor(service, 5000);
        m.restartService();

        mock.verify();
      });
    });

    suite('#restartIfTimedOut', function() {
      test('restarts when timed out', function() {
        const service = makeService();
        const mock = sinon.mock(service);
        mock.expects('stop').once();
        mock.expects('start').once();

        const m = new monitor.Monitor(service, 5000);
        sinon.stub(m, 'hasTimedOut').returns(true)

        m.restartIfTimedOut();

        mock.verify();
      });

      test('does not restart  when not timed out', function() {
        const service = makeService();
        const mock = sinon.mock(service);
        mock.expects('stop').never();
        mock.expects('start').never();

        const m = new monitor.Monitor(service, 5000);
        sinon.stub(m, 'hasTimedOut').returns(false)

        m.restartIfTimedOut();

        mock.verify();
      });
    });
  });
});
