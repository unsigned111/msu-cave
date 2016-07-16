/*jslint node: true, esversion: 6 */
'use strict';
var payloadValidator = require('payload-validator');

function requiredFields() {
  return [
    'delta', 'theta', 'loAlpha', 'hiAlpha', 'loBeta', 'hiBeta',
    'loGamma', 'midGamma'
  ];
}

function createExample(fields) {
  var source = {};
  fields.forEach(function(entry) { source[entry] = 0; });
  return source;
}

function validateBody(body) {
  // payload validator freaks out if you pass it a null, so
  // set the target to an empty object when none is provided
  let target = body || {};

  var fields = requiredFields();
  var example = createExample(fields);
  return payloadValidator.validator(target, example, fields, true);
}

module.exports.validate = validateBody;
