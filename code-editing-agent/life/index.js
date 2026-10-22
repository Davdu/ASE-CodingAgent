"use strict";
const { Life } = require("./life");

// Usage: node index.js [width] [height] [generations] [delayMs]
const [w = 40, h = 20, gens = 100, delay = 100] = process.argv.slice(2).map(Number);

const game = new Life(w, h, true);
// Seed with a glider plus random noise.
game.randomize(0.2);

const timer = setInterval(() => {
  console.clear();
  console.log(`Conway's Game of Life  gen=${game.generation}  pop=${game.population()}`);
  console.log(game.toString().replace(/#/g, "█").replace(/\./g, " "));
  if (game.generation >= gens || game.population() === 0) {
    clearInterval(timer);
    return;
  }
  game.step();
}, delay);
