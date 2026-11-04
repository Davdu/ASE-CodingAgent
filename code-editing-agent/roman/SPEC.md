# Roman numeral converter

Module `roman.js` exports two functions (CommonJS):

## `toRoman(n)`
1. Accepts an integer 1..3999 and returns its canonical Roman numeral string.
2. Uses subtractive notation (4 -> IV, 9 -> IX, 40 -> XL, 90 -> XC, 400 -> CD, 900 -> CM).
3. Throws `RangeError` for integers outside 1..3999.
4. Throws `TypeError` for non-integers (floats, NaN, strings, null, undefined).

## `fromRoman(s)`
1. Accepts a canonical, upper-case Roman numeral string and returns the integer.
2. Throws `SyntaxError` for non-canonical or invalid input (empty string, lower-case,
   unknown characters, "IIII", "VX", "IC", "IIX", ...).
3. Throws `TypeError` if the input is not a string.
4. Round trip: `fromRoman(toRoman(n)) === n` for every n in 1..3999.
