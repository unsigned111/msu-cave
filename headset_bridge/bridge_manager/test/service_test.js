/*jslint node: true, esversion: 6 */

const assert = require('chai').assert;
const sinon = require('sinon');

const service = require('../service');

suite('service', function() {
  suite('Service', function() {
    const config = {
      name: 'name',
      start: 'echo start',
      stop: 'echo stop',
    };

    suite('name', function() {
      it('sets name', function() {
        const s = new service.Service(config)
        assert('name', s.name);
      });
    });

    suite('#start', function() {
      it('runs a command with start command', function() {
        const s = new service.Service(config)
        spy = sinon.spy(s, 'runCmd');
        s.start();
        assert(spy.calledWith('echo start'));
      });
    });

    suite('#stop', function() {
      it('runs a command with stop command', function() {
        const s = new service.Service(config)
        spy = sinon.spy(s, 'runCmd');
        s.stop();
        assert(spy.calledWith('echo stop'));
      });

    });
  });
});


