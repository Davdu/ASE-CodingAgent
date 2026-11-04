const assert = require('node:assert');
const { test } = require('node:test');
const { toRoman, fromRoman } = require('./roman');

test('toRoman: basic values', () => {
  const cases = { 1: 'I', 2: 'II', 3: 'III', 4: 'IV', 5: 'V', 9: 'IX', 10: 'X',
    14: 'XIV', 40: 'XL', 49: 'XLIX', 90: 'XC', 400: 'CD', 900: 'CM',
    1994: 'MCMXCIV', 2024: 'MMXXIV', 3999: 'MMMCMXCIX' };
  for (const [n, r] of Object.entries(cases)) assert.strictEqual(toRoman(Number(n)), r);
});

test('toRoman: range errors', () => {
  for (const n of [0, -1, 4000, 100000]) assert.throws(() => toRoman(n), RangeError);
});

test('toRoman: type errors', () => {
  for (const v of [1.5, NaN, Infinity, '5', null, undefined]) assert.throws(() => toRoman(v), TypeError);
});

test('fromRoman: valid values', () => {
  const cases = { I: 1, IV: 4, IX: 9, XIV: 14, XLIX: 49, MCMXCIV: 1994, MMMCMXCIX: 3999 };
  for (const [r, n] of Object.entries(cases)) assert.strictEqual(fromRoman(r), n);
});

test('fromRoman: invalid strings', () => {
  for (const s of ['', 'abc', 'iv', 'IIII', 'VX', 'IC', 'IIX', 'MMMM', 'XXXX', 'VV', 'IXI'])
    assert.throws(() => fromRoman(s), SyntaxError, `should reject "${s}"`);
});

test('fromRoman: type errors', () => {
  for (const v of [5, null, undefined, {}]) assert.throws(() => fromRoman(v), TypeError);
});

test('round trip for 1..3999', () => {
  for (let n = 1; n <= 3999; n++) assert.strictEqual(fromRoman(toRoman(n)), n);
});
