'use strict';

const TABLE = [
  [1000, 'M'], [900, 'CM'], [500, 'D'], [400, 'CD'],
  [100, 'C'], [90, 'XC'], [50, 'L'], [40, 'XL'],
  [10, 'X'], [9, 'IX'], [5, 'V'], [4, 'IV'], [1, 'I'],
];

const VALID = /^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$/;

function toRoman(n) {
  if (typeof n !== 'number' || !Number.isInteger(n)) {
    throw new TypeError('Expected an integer');
  }
  if (n < 1 || n > 3999) throw new RangeError('Expected a value between 1 and 3999');
  let out = '';
  for (const [value, sym] of TABLE) {
    while (n >= value) {
      out += sym;
      n -= value;
    }
  }
  return out;
}

function fromRoman(s) {
  if (typeof s !== 'string') throw new TypeError('Expected a string');
  if (s === '' || !VALID.test(s)) throw new SyntaxError(`Invalid Roman numeral: "${s}"`);
  let total = 0;
  let rest = s;
  for (const [value, sym] of TABLE) {
    while (rest.startsWith(sym)) {
      total += value;
      rest = rest.slice(sym.length);
    }
  }
  return total;
}

module.exports = { toRoman, fromRoman };
