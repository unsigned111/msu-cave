/*jslint node: true, esversion: 6 */

var assert = require('chai').assert;

var validator = require('../validator');

suite('validator', function() {
  function goodBody() {
    return {
      delta: 1,
      theta: 2,
      loAlpha: 3,
      hiAlpha: 4,
      loBeta: 5,
      hiBeta: 6,
      loGamma: 7,
      midGamma: 8,
      timestamp: 9,
    };
  }

  suite('#validate()', function() {
    test('fails validation when body missing', function() {
      var result = validator.validate(null);
      assert.equal(result.success, false);
    });

    test('fails validation when body missing param', function() {
      var body = goodBody();
      delete body.midGamma;
      var result = validator.validate(body);
      assert.equal(result.success, false);
    });

    test('passes validation when body is valid', function() {
      var body = goodBody();
      var result = validator.validate(body);
      assert.equal(result.success, true);
    });
  });
});

