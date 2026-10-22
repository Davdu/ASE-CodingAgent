"use strict";

/** Conway's Game of Life on a finite grid (dead cells outside the edges). */
class Life {
  constructor(width, height, wrap = false) {
    this.width = width;
    this.height = height;
    this.wrap = wrap;
    this.grid = Array.from({ length: height }, () => new Array(width).fill(0));
    this.generation = 0;
  }

  static fromString(str, wrap = false) {
    const rows = str.trim().split("\n").map((r) => r.trim());
    const g = new Life(Math.max(...rows.map((r) => r.length)), rows.length, wrap);
    rows.forEach((row, y) => [...row].forEach((ch, x) => { g.grid[y][x] = ch === "#" || ch === "O" ? 1 : 0; }));
    return g;
  }

  set(x, y, alive = 1) { this.grid[y][x] = alive ? 1 : 0; return this; }

  get(x, y) {
    if (this.wrap) {
      x = (x + this.width) % this.width;
      y = (y + this.height) % this.height;
    } else if (x < 0 || y < 0 || x >= this.width || y >= this.height) {
      return 0;
    }
    return this.grid[y][x];
  }

  neighbors(x, y) {
    let n = 0;
    for (let dy = -1; dy <= 1; dy++)
      for (let dx = -1; dx <= 1; dx++)
        if (dx || dy) n += this.get(x + dx, y + dy);
    return n;
  }

  step() {
    const next = this.grid.map((row, y) =>
      row.map((cell, x) => {
        const n = this.neighbors(x, y);
        return n === 3 || (cell && n === 2) ? 1 : 0;
      })
    );
    this.grid = next;
    this.generation++;
    return this;
  }

  population() { return this.grid.flat().reduce((a, b) => a + b, 0); }

  randomize(density = 0.3, rng = Math.random) {
    this.grid = this.grid.map((row) => row.map(() => (rng() < density ? 1 : 0)));
    return this;
  }

  toString() {
    return this.grid.map((row) => row.map((c) => (c ? "#" : ".")).join("")).join("\n");
  }
}

module.exports = { Life };
