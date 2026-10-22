"use strict";
const assert = require("assert");
const { Life } = require("./life");

const tests = [];
const test = (name, fn) => tests.push([name, fn]);

test("block still life stays unchanged", () => {
  const g = Life.fromString(`
    ....
    .##.
    .##.
    ....`);
  const before = g.toString();
  g.step();
  assert.strictEqual(g.toString(), before);
});

test("blinker oscillates with period 2", () => {
  const g = Life.fromString(`
    .....
    ..#..
    ..#..
    ..#..
    .....`);
  const start = g.toString();
  g.step();
  assert.strictEqual(g.toString(), [".....", ".....", ".###.", ".....", "....."].join("\n"));
  g.step();
  assert.strictEqual(g.toString(), start);
  assert.strictEqual(g.generation, 2);
});

test("glider moves diagonally by (1,1) after 4 generations", () => {
  const g = Life.fromString(`
    .#....
    ..#...
    ###...
    ......
    ......
    ......`);
  for (let i = 0; i < 4; i++) g.step();
  assert.strictEqual(g.toString(), [
    "......",
    "..#...",
    "...#..",
    ".###..",
    "......",
    "......",
  ].join("\n"));
  assert.strictEqual(g.population(), 5);
});

test("lonely cell dies, overcrowded cell dies", () => {
  const g = new Life(3, 3).set(1, 1);
  g.step();
  assert.strictEqual(g.population(), 0);
  const full = Life.fromString("###\n###\n###");
  full.step();
  assert.strictEqual(full.get(1, 1), 0); // 8 neighbors -> dies
  assert.strictEqual(full.get(0, 0), 1); // 3 neighbors -> survives
});

test("wrapping glider keeps population of 5 forever", () => {
  // 4x4 torus is too small (glider collides with itself), use 8x8.
  const g = Life.fromString(".#......\n..#.....\n###.....\n........\n........\n........\n........\n........", true);
  for (let i = 0; i < 40; i++) g.step();
  assert.strictEqual(g.population(), 5);
});

let failed = 0;
for (const [name, fn] of tests) {
  try { fn(); console.log("PASS", name); }
  catch (e) { failed++; console.log("FAIL", name, "\n  ", e.message); }
}
console.log(`\n${tests.length - failed}/${tests.length} passed`);
process.exit(failed ? 1 : 0);
